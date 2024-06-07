package handler

import (
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/rest/httpx"
	"live/common/response"
	"live/servers/live_sever/live_api/internal/logic"
	"live/servers/live_sever/live_api/internal/svc"
	"live/servers/live_sever/live_api/internal/types"
	"net/http"
)

func peopleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PeopleRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		var upGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upGrader.Upgrade(w, r, nil)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		defer conn.Close()

		peopleLogic := logic.NewPeopleLogic(r.Context(), svcCtx)
		if err := peopleLogic.HandlePeople(req, conn); err != nil {
			response.Response(r, w, nil, err)
			return
		}
	}
}
