package logic

import (
	"context"
	"fmt"
	"live/models/live_models"
	"live/servers/live_sever/live_api/internal/svc"
	"live/servers/live_sever/live_api/internal/types"
	"live/utils/cache"
	"live/utils/stream"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.CreateRequest) (resp *types.CreateResponse, err error) {

	var liveModel live_models.LiveModel
	// 尝试找到一个已存在的 LiveModel
	roomNumber := fmt.Sprintf("%07d", req.AnchorID)

	newCache := cache.NewCache(l.svcCtx.Redis, l.svcCtx.DB)
	fetchFromDB := func() (interface{}, error) {
		_ = l.svcCtx.DB.Take(&liveModel, "anchor_id = ?", req.AnchorID).Error
		liveModel.RoomNumber = roomNumber
		liveModel.AnchorId = req.AnchorID
		liveModel.IsStart = true
		liveModel.Title = req.Title
		liveModel.Description = req.Description
		liveModel.RTMPAddress = "http://127.0.0.1:7001/live/" + roomNumber + ".flv"
		if req.Avatar != "" {
			liveModel.Avatar = req.Avatar
		} else {
			liveModel.Avatar = "../../../models/images/cover.png"
		}
		liveModel.AudienceCount = 0
		return liveModel, nil
	}
	expire := newCache.GetRandomExpire(24)
	data, err := newCache.GetOrSetCache(l.ctx, "roomInfo:"+strconv.Itoa(int(req.AnchorID)), fetchFromDB, live_models.LiveModel{}, expire)
	if err != nil {
		return nil, err
	}
	liveModel = data.(live_models.LiveModel)
	liveModel.Title = req.Title
	liveModel.Description = req.Description
	l.svcCtx.DB.Save(&liveModel)

	resp = &types.CreateResponse{
		RoomNumber: roomNumber,
		CreateAt:   time.Now().Format("2006-01-02 15:04:05"),
	}
	//生成RTMP推流地址
	channelKey := stream.GetChannelKey(resp.RoomNumber)
	rtmpEndpoint := "rtmp://localhost:1935/live/" + channelKey
	resp.RTMPEndpoint = rtmpEndpoint

	return resp, nil
}
