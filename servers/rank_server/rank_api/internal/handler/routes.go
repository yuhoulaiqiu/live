// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"live/servers/rank_server/rank_api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/anchor",
				Handler: GetAnchorFansRankHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/rank",
				Handler: GetRoomRankHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/rank/ws",
				Handler: WebSocketHandler(serverCtx),
			},
		},
	)
}