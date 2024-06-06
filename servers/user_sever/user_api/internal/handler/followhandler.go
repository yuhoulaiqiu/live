package handler

import (
	"live/common/response"
	"live/servers/user_sever/user_api/internal/logic"
	"live/servers/user_sever/user_api/internal/svc"
	"live/servers/user_sever/user_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func followHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FollowRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewFollowLogic(r.Context(), svcCtx)
		resp, err := l.Follow(&req)
		response.Response(r, w, resp, err)
	}
}
