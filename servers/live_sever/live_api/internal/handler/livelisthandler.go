package handler

import (
	"live/common/response"
	"live/servers/live_sever/live_api/internal/logic"
	"live/servers/live_sever/live_api/internal/svc"
	"net/http"
)

func liveListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		l := logic.NewLiveListLogic(r.Context(), svcCtx)
		resp, err := l.LiveList()
		response.Response(r, w, resp, err)
	}
}
