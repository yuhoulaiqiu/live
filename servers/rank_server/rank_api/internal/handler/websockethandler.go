package handler

import (
	"live/servers/rank_server/rank_api/internal/logic"
	"live/servers/rank_server/rank_api/internal/svc"
	"net/http"
)

func WebSocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wsService := logic.NewWebSocketService(r.Context(), svcCtx)
		wsService.HandleConnections(w, r)
	}
}
