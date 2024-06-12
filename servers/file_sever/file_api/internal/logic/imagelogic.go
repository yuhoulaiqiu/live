package logic

import (
	"context"

	"live/servers/file_sever/file_api/internal/svc"
	"live/servers/file_sever/file_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageLogic {
	return &ImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImageLogic) Image(req *types.ImageRequest) (resp *types.ImageResponse, err error) {
	resp = new(types.ImageResponse)

	return
}
