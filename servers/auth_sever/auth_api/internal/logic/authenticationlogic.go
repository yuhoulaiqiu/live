package logic

import (
	"context"
	"errors"
	"fmt"
	"live/servers/auth_sever/auth_api/internal/svc"
	"live/servers/auth_sever/auth_api/internal/types"
	"live/utils"
	"live/utils/jwts"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthenticationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthenticationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthenticationLogic {
	return &AuthenticationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthenticationLogic) Authentication(req *types.AuthenticationRequest) (resp *types.AuthenticationResponse, err error) {
	fmt.Println(req.Token, req.ValidPath)
	if utils.InListByRegex(l.svcCtx.Config.WhiteList, req.ValidPath) {
		logx.Infof("%s 在白名单中", req.ValidPath)
		return
	}

	if req.Token == "" {
		logx.Infof("token为空")
		err = errors.New("认证失败")
		return
	}
	claims, err := jwts.ParseToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		logx.Error("token解析失败")
		err = errors.New("认证失败")
		return
	}
	_, err = l.svcCtx.Redis.Get(fmt.Sprintf("logout_%s", req.Token)).Result()
	if err == nil {
		logx.Error("token已经登出")
		err = errors.New("认证失败")
		return
	}
	err = nil
	return &types.AuthenticationResponse{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}
