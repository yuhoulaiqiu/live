package logic

import (
	"context"
	"errors"
	"live/models/user_models"
	"live/servers/user_sever/user_api/internal/svc"
	"live/servers/user_sever/user_api/internal/types"
	"live/utils/cache"
	"strconv"

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
	//查缓存
	newCache := cache.NewCache(l.svcCtx.Redis, l.svcCtx.DB)
	fetchFromDB := func() (interface{}, error) {
		err = l.svcCtx.DB.Where("id = ?", req.UserID).First(&user).Error
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	expire := newCache.GetRandomExpire(24)
	data, err := newCache.GetOrSetCache(l.ctx, "userInfo:"+strconv.Itoa(int(req.UserID)), fetchFromDB, user_models.UserModel{}, expire)
	if err != nil {
		return nil, errors.New("获取用户信息失败")
	}
	user = data.(user_models.UserModel)

	resp = &types.UserInfoResponse{
		Username: user.UserName,
		Balance:  int(user.Balances),
		Nickname: user.NickName,
		Avatar:   user.Avatar,
		Fans:     user.Fans,
	}
	return
}
