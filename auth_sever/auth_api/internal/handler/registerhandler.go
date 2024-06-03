package handler

import (
	"live/auth_sever/auth_api/internal/logic"
	"live/auth_sever/auth_api/internal/svc"
	"live/auth_sever/auth_api/internal/types"
	"live/commen/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func registerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewRegisterLogic(r.Context(), svcCtx)
		err := l.Register(&req)
		response.Response(r, w, nil, err)
	}
}
