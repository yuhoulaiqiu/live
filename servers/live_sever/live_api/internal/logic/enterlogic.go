package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"live/models/live_models"
	"live/models/user_models"
	"live/servers/live_sever/live_api/internal/svc"
	"live/servers/live_sever/live_api/internal/types"
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
	err := l.svcCtx.DB.Where("room_number = ?", req.RoomNumber).First(&liveModel).Error
	if err != nil {
		return nil, errors.New("直播间不存在")
	}
	//改变用户所在直播间
	var userLiveModel user_models.UserModel
	err = l.svcCtx.DB.Where("id = ?", req.UserID).First(&userLiveModel).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if userLiveModel.InWhich == req.RoomNumber {
	} else {
		userLiveModel.InWhich = req.RoomNumber
		l.svcCtx.DB.Save(&userLiveModel)
	}
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
