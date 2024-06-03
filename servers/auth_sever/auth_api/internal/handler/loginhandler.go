package handler

import (
	"live/commen/response"
	"live/servers/auth_sever/auth_api/internal/logic"
	"live/servers/auth_sever/auth_api/internal/svc"
	"live/servers/auth_sever/auth_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func loginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		response.Response(r, w, resp, err)
	}
}
