package logic

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"live/models/live_models"
	"live/servers/live_sever/live_api/internal/svc"
	"live/servers/live_sever/live_api/internal/types"
	"live/utils/cache"
	"strconv"
)

type EndLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEndLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EndLogic {
	return &EndLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EndLogic) End(req *types.EndRequest) (resp *types.EndResponse, err error) {
	//判断是否有权限
	var liveModel live_models.LiveModel
	newCache := cache.NewCache(l.svcCtx.Redis, l.svcCtx.DB)
	fetchFromDB := func() (interface{}, error) {
		err = l.svcCtx.DB.Take(&liveModel, "anchor_id = ?", req.AnchorID).Error
		if err != nil {
			return nil, err
		}
		return liveModel, nil
	}
	expire := newCache.GetRandomExpire(24)
	data, err := newCache.GetOrSetCache(l.ctx, "roomInfo:"+strconv.Itoa(int(req.AnchorID)), fetchFromDB, live_models.LiveModel{}, expire)
	if err != nil {
		return nil, err
	}
	liveModel = data.(live_models.LiveModel)
	if liveModel.AnchorId != req.AnchorID {
		return nil, errors.New("没有权限")
	}
	//结束直播
	liveModel.IsStart = false
	//同步到redis
	dataBytes, _ := json.Marshal(liveModel)
	l.svcCtx.Redis.Set("roomInfo:"+strconv.Itoa(int(req.AnchorID)), dataBytes, expire)
	l.svcCtx.DB.Save(&liveModel)
	//删除redis中当日礼物信息
	l.svcCtx.Redis.Del("gift_ranking_" + strconv.Itoa(int(req.AnchorID)))
	return
}
