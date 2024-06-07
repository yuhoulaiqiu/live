package handler

import (
	"live/common/response"
	"live/servers/rank_server/rank_api/internal/logic"
	"live/servers/rank_server/rank_api/internal/svc"
	"live/servers/rank_server/rank_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAnchorFansRankHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAnchorFansRankRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetAnchorFansRankLogic(r.Context(), svcCtx)
		resp, err := l.GetAnchorFansRank(&req)
		response.Response(r, w, resp, err)
	}
}
