package handler

import (
	"live/common/response"
	"live/servers/live_sever/live_api/internal/logic"
	"live/servers/live_sever/live_api/internal/svc"
	"live/servers/live_sever/live_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func enterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EnterRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		l := logic.NewEnterLogic(r.Context(), svcCtx)
		resp, err := l.Enter(&req)
		response.Response(r, w, resp, err)
	}
}
