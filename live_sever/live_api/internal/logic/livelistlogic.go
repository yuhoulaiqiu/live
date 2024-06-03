package logic

import (
	"context"
	"live/models/live_models"

	"live/live_sever/live_api/internal/svc"
	"live/live_sever/live_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LiveListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLiveListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LiveListLogic {
	return &LiveListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LiveListLogic) LiveList() (resp *types.ListMessage, err error) {
	// 获取直播间信息
	var liveModel []live_models.LiveModel
	//查询所有开启的直播间
	err = l.svcCtx.DB.Where("is_start = ?", true).Find(&liveModel).Error
	if err != nil {
		return nil, err
	}
	var list []types.LiveMessage
	for _, v := range liveModel {
		list = append(list, types.LiveMessage{
			Title:       v.Title,
			Description: v.Description,
			OnlineUsers: v.AudienceCount,
			RoomNumber:  v.RoomNumber,
		})
	}
	resp = &types.ListMessage{
		LiveMessages: list,
	}

	return resp, nil
}
