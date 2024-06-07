package handler

import (
	"live/servers/rank_server/rank_api/internal/logic"
	"live/servers/rank_server/rank_api/internal/svc"
	"net/http"
)

func GiftWSHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wsService := logic.NewGiftWSLogic(r.Context(), svcCtx)
		wsService.HandleConnections(w, r)
	}
}
