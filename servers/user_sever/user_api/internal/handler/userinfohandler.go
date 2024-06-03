package handler

import (
	"github.com/zeromicro/go-zero/core/logx"
	"live/commen/response"
	"live/servers/user_sever/user_api/internal/logic"
	"live/servers/user_sever/user_api/internal/svc"
	"live/servers/user_sever/user_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func userInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserInfoRequest
		if err := httpx.Parse(r, &req); err != nil {
			logx.Errorf("httpx.Parse(r, &req) err: %v", err)
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.UserInfo(&req)
		response.Response(r, w, resp, err)
	}
}
