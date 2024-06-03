package logic

import (
	"context"
	"errors"
	"live/models/live_models"

	"live/live_sever/live_api/internal/svc"
	"live/live_sever/live_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EndLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEndLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EndLogic {
	return &EndLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EndLogic) End(req *types.EndRequest) (resp *types.EndResponse, err error) {
	//判断是否有权限
	var liveModel live_models.LiveModel
	err = l.svcCtx.DB.Where("room_number = ?", req.RoomNumber).First(&liveModel).Error
	if err != nil {
		return nil, errors.New("直播间不存在")
	}
	if liveModel.AnchorId != req.AnchorID {
		return nil, errors.New("没有权限")
	}
	//结束直播
	liveModel.IsStart = false
	l.svcCtx.DB.Save(&liveModel)
	return
}
