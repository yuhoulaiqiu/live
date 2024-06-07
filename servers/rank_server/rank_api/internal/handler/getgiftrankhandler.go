package handler

import (
	"live/common/response"
	"live/servers/rank_server/rank_api/internal/logic"
	"live/servers/rank_server/rank_api/internal/svc"
	"live/servers/rank_server/rank_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetGiftRankHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetGiftRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetGiftRankLogic(r.Context(), svcCtx)
		resp, err := l.GetGiftRank(&req)
		response.Response(r, w, resp, err)
	}
}
