package logic

import (
	"context"
	"errors"
	"live/models/user_models"

	"live/user_sever/user_api/internal/svc"
	"live/user_sever/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	var user user_models.UserModel
	err = l.svcCtx.DB.Where("id = ?", req.UserID).First(&user).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	resp = &types.UserInfoResponse{
		Nickname: user.NickName,
		Avatar:   user.Avatar,
		Fans:     user.Fans,
	}
	return
}
