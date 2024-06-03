package logic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"live/models/live_models"
	"live/servers/live_sever/live_api/internal/svc"
	"live/servers/live_sever/live_api/internal/types"
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
	number := "0000000" + strconv.Itoa(int(req.AnchorID))
	roomNumber := number[len(number)-7:]
	err = l.svcCtx.DB.Where("anchor_id = ?", req.AnchorID).First(&liveModel).Error
	if err != nil {
		// 如果没有找到 LiveModel，创建一个新的
		if errors.Is(err, gorm.ErrRecordNotFound) {
			liveModel = live_models.LiveModel{
				Title:         req.Title,
				AnchorId:      req.AnchorID,
				Description:   req.Description,
				RoomNumber:    roomNumber,
				IsStart:       true,
				Avatar:        req.Avatar,
				RTMPAddress:   "http://127.0.0.1:7001/live/" + liveModel.RoomNumber + ".flv",
				AudienceCount: 0,
			}
			err = l.svcCtx.DB.Create(&liveModel).Error
			if err != nil {
				return nil, errors.New("创建直播间失败")
			}
		} else {
			// 如果错误是其他的，返回它
			return nil, err
		}
	} else {
		liveModel.IsStart = true
		liveModel.Title = req.Title
		liveModel.Description = req.Description
		liveModel.RTMPAddress = "http://127.0.0.1:7001/live/" + liveModel.RoomNumber + ".flv"
		if req.Avatar != "" {
			liveModel.Avatar = req.Avatar
		} else {
			liveModel.Avatar = "../../../images/Live broadcast.png"
		}
		liveModel.AudienceCount = 0

		// 保存 LiveModel
		l.svcCtx.DB.Save(&liveModel)
	}
	resp = &types.CreateResponse{
		RoomNumber: roomNumber,
		CreateAt:   time.Now().Format("2006-01-02 15:04:05"),
	}
	//生成RTMP推流地址
	channelKey := stream.GetChannelKey(resp.RoomNumber)
	rtmpEndpoint := "rtmp://localhost:1935/live/" + channelKey
	resp.RTMPEndpoint = rtmpEndpoint
	//开始推流
	//这部分通常是在主播的客户端应用中完成的，后端只需要提供RTMP推流地址
	return resp, nil
}
