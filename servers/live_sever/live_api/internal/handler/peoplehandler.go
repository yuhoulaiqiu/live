package handler

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"live/common/response"
	"live/models/live_models"
	"live/servers/live_sever/live_api/internal/svc"
	"live/servers/live_sever/live_api/internal/types"
	"net/http"
	"sync"
	"time"
)

// 观众列表
var audienceList = make(map[string]map[*websocket.Conn]bool)
var lock = sync.RWMutex{}

// 存储每个房间的上下文和取消函数
var roomContexts = make(map[string]context.CancelFunc)

func peopleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PeopleRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		// 判断房间号是否存在
		var liveModel live_models.LiveModel
		err := svcCtx.DB.Where("room_number = ?", req.RoomNumber).First(&liveModel).Error
		if err != nil {
			response.Response(r, w, nil, errors.New("房间号不存在"))
			return
		}
		var upGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upGrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}

		// 将新的用户添加到观众列表中
		lock.Lock()
		if _, ok := audienceList[req.RoomNumber]; !ok {
			audienceList[req.RoomNumber] = make(map[*websocket.Conn]bool)
		}
		audienceList[req.RoomNumber][conn] = true
		lock.Unlock()

		// 更新 Redis 中的实时人数
		svcCtx.Redis.ZAdd("room_ranking", redis.Z{
			Score:  float64(len(audienceList[req.RoomNumber])),
			Member: req.RoomNumber,
		})

		// 检查是否已经有一个协程在为这个房间服务
		if _, ok := roomContexts[req.RoomNumber]; !ok {
			// 创建一个可以被取消的上下文
			ctx, cancel := context.WithCancel(context.Background())
			// 存储上下文和取消函数
			roomContexts[req.RoomNumber] = cancel
			// 启动定时任务，每1秒同步一次数据到 MySQL
			go syncAudienceCountToDB(ctx, svcCtx)
		}

		// 发送实时人数
		for conn := range audienceList[req.RoomNumber] {
			conn.WriteJSON(len(audienceList[req.RoomNumber]))
		}
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				// 用户断开连接，将其从观众列表中移除
				lock.Lock()
				delete(audienceList[req.RoomNumber], conn)
				lock.Unlock()
				// 更新 Redis 中的实时人数
				svcCtx.Redis.ZAdd("room_ranking", redis.Z{
					Score:  float64(len(audienceList[req.RoomNumber])),
					Member: req.RoomNumber,
				})
				// 发送更新后的实时人数
				for conn := range audienceList[req.RoomNumber] {
					conn.WriteJSON(len(audienceList[req.RoomNumber]))
				}
				break
			}
		}
		defer conn.Close()
	}
}

// 定时任务，每1秒同步一次数据到 MySQL
func syncAudienceCountToDB(ctx context.Context, svcCtx *svc.ServiceContext) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for roomNumber, audience := range audienceList {
				// 同步到数据库
				var liveModel live_models.LiveModel
				err := svcCtx.DB.Where("room_number = ?", roomNumber).First(&liveModel).Error
				if err != nil {
					logx.Error(err)
					break
				}
				// 检查直播间是否已经关闭
				if liveModel.IsStart == false {
					// 直播间已经关闭，取消上下文，停止 goroutine
					if cancel, ok := roomContexts[roomNumber]; ok {
						cancel()
						delete(roomContexts, roomNumber)
					}
					break
				}
				liveModel.AudienceCount = len(audience)
				svcCtx.DB.Save(&liveModel)
			}
		case <-ctx.Done():
			// 上下文被取消，停止 goroutine
			// 在停止前，最后一次同步数据到数据库
			for roomNumber, audience := range audienceList {
				var liveModel live_models.LiveModel
				err := svcCtx.DB.Where("room_number = ?", roomNumber).First(&liveModel).Error
				if err != nil {
					logx.Error(err)
					break
				}
				liveModel.AudienceCount = len(audience)
				svcCtx.DB.Save(&liveModel)
			}
			logx.Infof("协程已经停止")
			return
		}
	}
}
