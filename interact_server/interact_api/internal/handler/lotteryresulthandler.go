package handler

import (
	"live/commen/response"
	"live/interact_server/interact_api/internal/logic"
	"live/interact_server/interact_api/internal/svc"
	"live/interact_server/interact_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func lotteryResultHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LotteryResultRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewLotteryResultLogic(r.Context(), svcCtx)
		resp, err := l.LotteryResult(&req)
		response.Response(r, w, resp, err)
	}
}
