package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"live/models/live_models"
	"live/servers/live_sever/live_api/internal/svc"
	"live/servers/live_sever/live_api/internal/types"
	"live/utils/cache"
)

type EnterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnterLogic {
	return &EnterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnterLogic) Enter(req *types.EnterRequest) (*types.EnterResponse, error) {
	// 获取直播间信息
	var liveModel live_models.LiveModel
	//先查缓存
	newCache := cache.NewCache(l.svcCtx.Redis, l.svcCtx.DB)
	fetchFromDB := func() (interface{}, error) {
		err := l.svcCtx.DB.Where("room_number = ?", req.RoomNumber).First(&liveModel).Error
		if err != nil {
			return nil, err
		}
		return liveModel, nil
	}
	expire := newCache.GetRandomExpire(24)
	data, err := newCache.GetOrSetCache(l.ctx, "roomInfo:"+req.RoomNumber, fetchFromDB, live_models.LiveModel{}, expire)
	if err != nil {
		logx.Errorf("获取直播间信息失败: %v", err)
		return nil, errors.New("获取直播间信息失败")
	}
	liveModel = data.(live_models.LiveModel)

	// 获取直播流地址
	RTMPAddress := "http://127.0.0.1:7001/live/" + req.RoomNumber + ".flv"
	resp := &types.EnterResponse{
		LiveMessage: types.LiveMessage{
			Title:       liveModel.Title,
			Description: liveModel.Description,
			OnlineUsers: liveModel.AudienceCount,
			RoomNumber:  liveModel.RoomNumber,
		}, RTMPAddress: RTMPAddress,
	}
	l.svcCtx.DB.Save(&liveModel)
	// 返回响应
	return resp, nil
}
