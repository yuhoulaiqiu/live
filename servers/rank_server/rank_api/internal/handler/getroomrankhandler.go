package handler

import (
	"live/common/response"
	"live/servers/rank_server/rank_api/internal/logic"
	"live/servers/rank_server/rank_api/internal/svc"
	"live/servers/rank_server/rank_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRoomRankHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRoomRankRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetRoomRankLogic(r.Context(), svcCtx)
		resp, err := l.GetRoomRank(&req)
		response.Response(r, w, resp, err)
	}
}
