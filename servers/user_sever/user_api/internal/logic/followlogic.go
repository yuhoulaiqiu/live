package logic

import (
	"context"
	"errors"
	"live/models/user_models"

	"live/servers/user_sever/user_api/internal/svc"
	"live/servers/user_sever/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowLogic) Follow(req *types.FollowRequest) (resp *types.FollowResponse, err error) {
	var follow user_models.FansModel
	var anchor user_models.UserModel
	// 查询是否已经关注
	err = l.svcCtx.DB.Where("user_id = ? and anchor_id = ?", req.UserID, req.AnchorId).First(&follow).Error
	if err != nil {
		// 创建信息
		follow = user_models.FansModel{
			AnchorId: uint32(req.AnchorId),
			UserId:   uint32(req.UserID),
			IsFans:   true,
			Level:    1,
		}
		err = l.svcCtx.DB.Create(&follow).Error
		if err != nil {
			return nil, errors.New("关注失败")
		}
	} else {
		// 修改关注状态
		follow.IsFans = !follow.IsFans
		err = l.svcCtx.DB.Save(&follow).Error
		if err != nil {
			return nil, errors.New("关注失败")
		}
	}
	// 查询主播信息
	err = l.svcCtx.DB.Where("id = ?", req.AnchorId).First(&anchor).Error
	if follow.IsFans {
		// 关注
		anchor.Fans++
	} else {
		// 取消关注
		anchor.Fans--
	}
	err = l.svcCtx.DB.Save(&anchor).Error
	return
}
