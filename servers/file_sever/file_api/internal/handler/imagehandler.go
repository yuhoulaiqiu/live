package handler

import (
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"live/common/response"
	"live/servers/file_sever/file_api/internal/logic"
	"live/servers/file_sever/file_api/internal/svc"
	"live/servers/file_sever/file_api/internal/types"
	"live/utils"
	"live/utils/random"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

var imageTypeList = []string{"avatar", "cover"}

func ImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ImageRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		// 获取上传的文件
		file, fileHead, err := r.FormFile("image")
		if err != nil {
			if errors.Is(err, io.EOF) {
				// 如果错误是 EOF，返回一个更友好的错误消息
				response.Response(r, w, nil, errors.New("上传的文件在传输过程中被截断，请重试"))
			} else {
				// 其他错误，直接返回
				logx.Error(err)
				response.Response(r, w, nil, err)
			}
			return
		}
		// 获取图片类型
		imageType := r.FormValue("imageType")
		if imageType == "" {
			response.Response(r, w, nil, errors.New("imageType不能为空"))
			return
		}
		if !utils.InList(imageTypeList, imageType) {
			response.Response(r, w, nil, errors.New("imageType不合法"))
			return
		}
		// 检查文件大小
		mSize := float64(fileHead.Size) / float64(1024) / float64(1024)
		if svcCtx.Config.FileSize < mSize {
			response.Response(r, w, nil, fmt.Errorf("文件大小不能超过%.2fMB", svcCtx.Config.FileSize))
			return
		}
		// 检查文件后缀是否在白名单中
		suffix := path.Ext(fileHead.Filename)
		if !utils.InList(svcCtx.Config.WhiteList, suffix) {
			response.Response(r, w, nil, errors.New("文件类型不允许"))
			return
		}
		// 检查文件是否重名
		dirPath := path.Join(svcCtx.Config.UploadDir, imageType)
		dir, err := os.ReadDir(dirPath)
		if err != nil {
			os.MkdirAll(dirPath, 0666)
		}
		// 读取文件数据
		imageData, _ := io.ReadAll(file)
		fileName := fileHead.Filename
		filePath := path.Join(svcCtx.Config.UploadDir, imageType, fileName)
		logx.Infof("file path: %s", filePath)
		// 创建逻辑处理器
		l := logic.NewImageLogic(r.Context(), svcCtx)
		resp, _ := l.Image(&req)
		resp.Url = filePath
		// 如果文件已存在，检查是否重复
		if InDir(dir, fileName) {
			byteData, _ := os.ReadFile(filePath)
			oldFileHash := utils.MD5(imageData)
			newFileHash := utils.MD5(byteData)
			// 如果文件重复，不上传
			if oldFileHash == newFileHash {
				response.Response(r, w, resp, nil)
				return
			}
			// 如果文件不重复，修改文件名
			xxxx := random.RandomStr(4)
			fileName = strings.TrimSuffix(fileName, suffix) + "_" + xxxx + suffix
			filePath = path.Join(svcCtx.Config.UploadDir, imageType, fileName)
		}
		// 写入文件
		err = os.WriteFile(filePath, imageData, 0666)
		if err != nil {
			response.Response(r, w, nil, errors.New("上传失败"))
			return
		}

		// 返回响应
		resp, err = l.Image(&req)
		resp.Url = "/" + filePath
		response.Response(r, w, resp, err)
	}
}

func InDir(dir []os.DirEntry, name string) bool {
	for _, v := range dir {
		if v.Name() == name {
			return true
		}
	}
	return false
}
