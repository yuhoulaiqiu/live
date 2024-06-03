package handler

import (
	"live/auth_sever/auth_api/internal/logic"
	"live/auth_sever/auth_api/internal/svc"
	"live/commen/response"
	"net/http"
)

func logoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewLogoutLogic(r.Context(), svcCtx)
		token := r.Header.Get("token")
		resp, err := l.Logout(token)
		response.Response(r, w, resp, err)
	}
}
