package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"live/models/user_models"
	"live/servers/auth_sever/auth_api/internal/svc"
	"live/servers/auth_sever/auth_api/internal/types"
	"live/utils/cache"
	"live/utils/jwts"
	"live/utils/pwd"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	var user user_models.UserModel
	newCache := cache.NewCache(l.svcCtx.Redis, l.svcCtx.DB)
	fetchFromDB := func() (interface{}, error) {
		err = l.svcCtx.DB.Take(&user, "id = ?", req.UserName).Error
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	// 从缓存中获取用户信息，过期时间为24小时+一个随机数
	expire := newCache.GetRandomExpire(24)
	data, err := newCache.GetOrSetCache(l.ctx, "userInfo:"+req.UserName, fetchFromDB, user_models.UserModel{}, expire)
	if err != nil {
		logx.Errorf("获取用户信息失败:%v", err)
		return nil, errors.New("用户名或密码错误")
	}
	user = data.(user_models.UserModel)
	if !pwd.CheckPwd(user.Pwd, req.Password) {
		err = errors.New("用户名或密码错误")
		return nil, err
	}
	token, err := jwts.GenToken(jwts.JwtPayload{
		UserID:   user.ID,
		Nickname: user.NickName,
		Role:     int(user.Role),
	}, l.svcCtx.Config.Auth.AccessSecret, int(l.svcCtx.Config.Auth.AccessExpire))
	if err != nil {
		logx.Error(err)
		err = errors.New("生成token失败")
		return nil, err
	}

	return &types.LoginResponse{Token: token}, nil
}
