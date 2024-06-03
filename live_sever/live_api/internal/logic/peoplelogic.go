package logic

import (
	"context"

	"live/live_sever/live_api/internal/svc"
	"live/live_sever/live_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PeopleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPeopleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PeopleLogic {
	return &PeopleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PeopleLogic) People(req *types.PeopleRequest) (resp *types.PeopleRespones, err error) {
	// todo: add your logic here and delete this line

	return
}
