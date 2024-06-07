package handler

import (
	"live/servers/live_sever/live_api/internal/logic"
	"live/servers/live_sever/live_api/internal/svc"
	"net/http"
)

func peopleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wsService := logic.NewPeopleLogic(r.Context(), svcCtx)
		wsService.HandleConnections(w, r)
	}
}
