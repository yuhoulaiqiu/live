package logic

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"live/common/response"
	"live/servers/rank_server/rank_api/internal/svc"
	"live/servers/rank_server/rank_api/internal/types"
	"net/http"
	"time"
)

type WebSocketService struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	clients   map[*websocket.Conn]bool
	broadcast chan []types.RoomItem
}

func NewWebSocketService(ctx context.Context, svcCtx *svc.ServiceContext) *WebSocketService {
	return &WebSocketService{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []types.RoomItem),
	}
}

func (s *WebSocketService) HandleConnections(w http.ResponseWriter, r *http.Request) {
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
	s.clients[conn] = true

	go s.readMessages(conn)
	go s.writeMessages(conn)
	go s.broadcastRankUpdates()

	// Add this line to keep the connection open
	select {}
}
func (s *WebSocketService) readMessages(conn *websocket.Conn) {
	defer conn.Close()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(s.clients, conn)
			break
		}
	}
}

func (s *WebSocketService) writeMessages(conn *websocket.Conn) {
	for {
		rank := <-s.broadcast
		for client := range s.clients {
			err := client.WriteJSON(rank)
			if err != nil {
				logx.Error("websocket error:", err)
				client.Close()
				delete(s.clients, client)
			}
		}
	}
}

func (s *WebSocketService) broadcastRankUpdates() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		logic := NewGetRoomRankLogic(s.ctx, s.svcCtx)
		req := &types.GetRoomRankRequest{TopN: 10} // 假设获取前10名
		resp, err := logic.GetRoomRank(req)
		if err != nil {
			logx.Error(err)
			continue
		}
		s.broadcast <- resp.Ranks
	}
}
