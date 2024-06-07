package logic

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"live/common/response"
	"live/servers/rank_server/rank_api/internal/svc"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

//var giftRankingList = make(map[string]map[*websocket.Conn]bool)
//var giftLock = sync.RWMutex{}
//var giftRoomContexts = make(map[string]context.CancelFunc)

type Gift struct {
	UserID int     `json:"user_id"`
	Score  float64 `json:"score"`
}
type GiftRankingItem struct {
	Gifts []Gift `json:"gifts"`
}

type GiftWSLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	rooms     map[string]map[*websocket.Conn]bool
	broadcast map[string]chan GiftRankingItem
	giftLock  sync.RWMutex
}

func NewGiftWSLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GiftWSLogic {
	return &GiftWSLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		rooms:     make(map[string]map[*websocket.Conn]bool),
		broadcast: make(map[string]chan GiftRankingItem),
		giftLock:  sync.RWMutex{},
	}
}

func (s *GiftWSLogic) HandleConnections(w http.ResponseWriter, r *http.Request) {
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

	if _, ok := s.rooms[roomNumber]; !ok {
		s.rooms[roomNumber] = make(map[*websocket.Conn]bool)
		s.broadcast[roomNumber] = make(chan GiftRankingItem)
	}

	s.rooms[roomNumber][conn] = true

	go s.readMessages(conn, roomNumber)
	go s.writeMessages(conn, roomNumber)
	go s.broadcastRankUpdates(roomNumber)

	select {}
}

func (s *GiftWSLogic) readMessages(conn *websocket.Conn, roomNumber string) {
	defer conn.Close()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// 用户断开连接，将其从客户端列表中移除
			s.giftLock.Lock()
			delete(s.rooms[roomNumber], conn)
			s.giftLock.Unlock()
			break
		}
	}
}

func (s *GiftWSLogic) writeMessages(conn *websocket.Conn, roomNumber string) {
	for {
		select {
		case message := <-s.broadcast[roomNumber]:
			err := conn.WriteJSON(message)
			if err != nil {
				logx.Error("websocket error:", err)
				conn.Close()
				s.giftLock.Lock()
				delete(s.rooms[roomNumber], conn)
				s.giftLock.Unlock()
				return
			}
		}
	}
}

func (s *GiftWSLogic) broadcastRankUpdates(roomNumber string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 获取礼物排行榜
			roomKey := "gift_ranking_" + roomNumber[6:]
			result, err := s.svcCtx.Redis.ZRevRangeWithScores(roomKey, 0, 9).Result()
			if err != nil {
				logx.Error(err)
				continue
			}

			var giftRankings []Gift
			for _, z := range result {
				userID, _ := strconv.Atoi(z.Member.(string))
				giftRankings = append(giftRankings, Gift{
					UserID: userID,
					Score:  z.Score,
				})
			}

			message := GiftRankingItem{Gifts: giftRankings}

			// 发送实时礼物排行榜
			s.giftLock.RLock()
			for conn := range s.rooms[roomNumber] {
				err := conn.WriteJSON(message)
				if err != nil {
					logx.Error("websocket error:", err)
					conn.Close()
					s.giftLock.Lock()
					delete(s.rooms[roomNumber], conn)
					s.giftLock.Unlock()
				}
			}
			s.giftLock.RUnlock()
		}
	}
}
