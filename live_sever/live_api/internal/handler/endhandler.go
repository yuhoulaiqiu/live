package handler

import (
	"live/commen/response"
	"live/live_sever/live_api/internal/logic"
	"live/live_sever/live_api/internal/svc"
	"live/live_sever/live_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func endHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EndRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		l := logic.NewEndLogic(r.Context(), svcCtx)
		resp, err := l.End(&req)
		response.Response(r, w, resp, err)
	}
}
