package handler

import (
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"live/common/response"
	"live/servers/file_sever/file_api/internal/svc"
	"live/servers/file_sever/file_api/internal/types"
	"net/http"
	"os"
)

func ImageShowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ImageShowRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		filePath := "D:/Golang_backend/live/models/images/" + req.ImageType + "/" + req.ImageName
		byteData, err := os.ReadFile(filePath)
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, errors.New("图片不存在"))
			return
		}
		w.Write(byteData)
	}
}
