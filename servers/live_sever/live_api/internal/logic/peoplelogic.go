package logic

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"live/common/response"
	"live/models/live_models"
	"live/servers/live_sever/live_api/internal/svc"
	"net/http"
	"sync"
	"time"
)

// 直播间全局化，一个直播间对应一个ws连接
var rooms = make(map[string]map[*websocket.Conn]bool)
var Lock = sync.RWMutex{}

type PeopleLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	broadcast map[string]chan int
}

func NewPeopleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PeopleLogic {
	return &PeopleLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		broadcast: make(map[string]chan int),
	}
}
func (p *PeopleLogic) HandleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logx.Error(err)
		response.Response(r, w, nil, errors.New("websocket error"))
		return
	}

	defer conn.Close()

	roomNumber := r.URL.Query().Get("roomNumber")
	if roomNumber == "" {
		response.Response(r, w, nil, errors.New("roomNumber is required"))
		return
	}
	//判断房间号是否存在
	var liveModel live_models.LiveModel
	err = p.svcCtx.DB.Where("room_number = ?", roomNumber).First(&liveModel).Error
	if err != nil {
		response.Response(r, w, nil, errors.New("房间号不存在"))
		return
	}

	Lock.Lock()
	if _, ok := rooms[roomNumber]; !ok {
		rooms[roomNumber] = make(map[*websocket.Conn]bool)
		p.broadcast[roomNumber] = make(chan int)
	}
	rooms[roomNumber][conn] = true
	Lock.Unlock()

	go p.readMessages(conn, roomNumber)
	go p.writeMessages(conn, roomNumber)
	go p.broadcastRankUpdates(roomNumber)

	select {}
}

func (p *PeopleLogic) readMessages(conn *websocket.Conn, roomNumber string) {
	defer conn.Close()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// 用户断开连接，将其从客户端列表中移除
			Lock.Lock()
			delete(rooms[roomNumber], conn)
			Lock.Unlock()
			break
		}
	}
}

func (p *PeopleLogic) writeMessages(conn *websocket.Conn, roomNumber string) {
	for {
		message := <-p.broadcast[roomNumber]
		for client := range rooms[roomNumber] {
			err := client.WriteJSON(message)
			if err != nil {
				logx.Error("websocket error:", err)
				client.Close()
				Lock.Lock()
				delete(rooms[roomNumber], client)
				Lock.Unlock()
			}
		}
	}
}

func (p *PeopleLogic) broadcastRankUpdates(roomNumber string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 更新 Redis 中的实时人数
			p.svcCtx.Redis.ZAdd("room_ranking", redis.Z{
				Score:  float64(len(rooms[roomNumber])),
				Member: roomNumber,
			})
			// 发送实时人数
			p.broadcast[roomNumber] <- len(rooms[roomNumber])
		}
	}
}
