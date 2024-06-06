package handler

import (
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"live/common/response"
	"live/models/interact_models"
	"live/models/user_models"
	"live/servers/interact_server/interact_api/internal/svc"
	"live/servers/interact_server/interact_api/internal/types"
	"net/http"
	"sync"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 直播间列表
var roomList = make(map[string]map[*websocket.Conn]bool)
var lock = sync.RWMutex{}

func chatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatRequest
		if err := httpx.Parse(r, &req); err != nil {
			logx.Errorf("httpx.Parse(r, &req) err:%v", err)
			response.Response(r, w, nil, err)
			return
		}

		var upGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upGrader.Upgrade(w, r, nil)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 将新的用户添加到直播间列表中
		lock.Lock()
		if _, ok := roomList[req.RoomNumber]; !ok {
			roomList[req.RoomNumber] = make(map[*websocket.Conn]bool)
		}
		roomList[req.RoomNumber][conn] = true
		lock.Unlock()
		//获取用户nickname
		var user user_models.UserModel
		svcCtx.DB.Where("id = ?", req.UserId).First(&user)

		conn.WriteMessage(websocket.TextMessage, []byte("系统消息："+user.NickName+" 欢迎来到直播间"))
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				// 用户断开连接，将其从直播间列表中移除
				lock.Lock()
				delete(roomList[req.RoomNumber], conn)
				lock.Unlock()
				break
			}
			//入库
			var chat interact_models.ChatModel
			chat.RoomNumber = req.RoomNumber
			chat.SendUserId = req.UserId
			chat.Msg = string(message)
			svcCtx.DB.Create(&chat)
			// 将消息添加上用户昵称
			message = []byte(user.NickName + "：" + string(message))
			// 将消息发送给直播间的所有其他用户
			for userConn := range roomList[req.RoomNumber] {
				if userConn != conn {
					userConn.WriteMessage(websocket.TextMessage, message)
				}
			}
		}
		defer conn.Close()
	}
}
