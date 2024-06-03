// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"live/servers/live_sever/live_api/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/live/create",
				Handler: createHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/live/end",
				Handler: endHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/live/enter",
				Handler: enterHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/live/list",
				Handler: liveListHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/live/ws/people",
				Handler: peopleHandler(serverCtx),
			},
		},
	)
}