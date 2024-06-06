package handler

import (
	"live/common/response"
	"live/servers/interact_server/interact_api/internal/logic"
	"live/servers/interact_server/interact_api/internal/svc"
	"live/servers/interact_server/interact_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func participateLotteryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ParticipateLotteryRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewParticipateLotteryLogic(r.Context(), svcCtx)
		resp, err := l.ParticipateLottery(&req)
		response.Response(r, w, resp, err)
	}
}
