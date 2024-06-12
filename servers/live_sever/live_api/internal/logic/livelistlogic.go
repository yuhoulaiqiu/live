package logic

import (
	"context"
	"errors"
	"live/models/live_models"
	"live/servers/live_sever/live_api/internal/svc"
	"live/servers/live_sever/live_api/internal/types"
	"live/utils/cache"

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
	newCache := cache.NewCache(l.svcCtx.Redis, l.svcCtx.DB)
	fetchFromDB := func() (interface{}, error) {
		err = l.svcCtx.DB.Where("is_start = ?", true).Find(&liveModel).Error
		if err != nil {
			return nil, err
		}
		return liveModel, nil
	}
	expire := newCache.GetRandomExpire(24)
	data, err := newCache.GetOrSetCache(l.ctx, "liveList", fetchFromDB, []live_models.LiveModel{}, expire)
	if err != nil {
		return nil, errors.New("获取直播列表失败")
	}
	liveModel = data.([]live_models.LiveModel)
	var list []types.LiveMessage
	for _, v := range liveModel {
		//从redis中获取直播间的观众人数
		onlineUsers, err := l.svcCtx.Redis.ZScore("room_ranking", v.RoomNumber).Result()
		if err != nil {
			logx.Error("从 Redis 获取在线用户失败:", err)
			return nil, errors.New("获取直播间观众人数失败")
		}
		list = append(list, types.LiveMessage{
			Title:       v.Title,
			Description: v.Description,
			OnlineUsers: int(onlineUsers),
			RoomNumber:  v.RoomNumber,
		})
	}
	resp = &types.ListMessage{
		LiveMessages: list,
	}

	return resp, nil
}
