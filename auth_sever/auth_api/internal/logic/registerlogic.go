package logic

import (
	"context"
	"errors"
	"live/models/user_models"
	"live/utils/pwd"

	"live/auth_sever/auth_api/internal/svc"
	"live/auth_sever/auth_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.LoginRequest) error {
	var user user_models.UserModel
	err := l.svcCtx.DB.Where("user_name = ?", req.UserName).First(&user).Error
	if err == nil {
		return errors.New("用户名已存在")
	}
	user = user_models.UserModel{
		UserName: req.UserName,
		NickName: "观众",
		Pwd:      pwd.HashPwd(req.Password),
		Role:     2,
		Avatar:   "../../../images/avatar.png",
		Fans:     0,
	}
	err = l.svcCtx.DB.Create(&user).Error
	if err != nil {
		return errors.New("注册失败")
	}

	return nil
}
