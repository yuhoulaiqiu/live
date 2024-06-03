package logic

import (
	"context"
	"encoding/json"
	"errors"
	"live/models/interact_models"
	"live/servers/interact_server/interact_api/internal/svc"
	"live/servers/interact_server/interact_api/internal/types"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type ParticipateLotteryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewParticipateLotteryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParticipateLotteryLogic {
	return &ParticipateLotteryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ParticipateLotteryLogic) ParticipateLottery(req *types.ParticipateLotteryRequest) (resp *types.ParticipateLotteryResponse, err error) {
	// 查找抽奖信息
	var lottery interact_models.LotteryModel
	if err := l.svcCtx.DB.Where("id = ? AND is_completed = ?", req.LotteryId, false).First(&lottery).Error; err != nil {
		logx.Error("抽奖不存在或已结束")
		return nil, errors.New("抽奖不存在或已结束")
	}
	//防止用户恶意调用接口，多次参与抽奖
	var count int64
	//查redis
	count, err = l.svcCtx.Redis.HLen("lottery:" + strconv.Itoa(int(req.LotteryId)) + ":participants").Result()
	if err != nil {
		logx.Error("查询参与用户失败")
		return nil, err
	}
	if count > 0 {
		logx.Error("用户已参与抽奖")
		return nil, errors.New("用户已参与抽奖")

	}
	//// 创建参与信息
	//participation := &interact_models.LotteryParticipationModel{
	//	LotteryId:  req.LotteryId,
	//	UserId:     req.UserId,
	//	MethodType: req.MethodType,
	//}
	//
	//// 保存参与信息到数据库
	//if err := l.svcCtx.DB.Create(&participation).Error; err != nil {
	//	logx.Error("参与抽奖失败")
	//	return nil, err
	//}
	// 创建参与信息
	participation := &interact_models.LotteryParticipationModel{
		LotteryId:  req.LotteryId,
		UserId:     req.UserId,
		MethodType: req.MethodType,
	}

	// 将参与信息转换为 JSON
	participationJson, err := json.Marshal(participation)
	if err != nil {
		logx.Error("参与抽奖失败")
		return nil, err
	}

	// 保存参与信息到 Redis
	err = l.svcCtx.Redis.HSet("lottery:"+strconv.Itoa(int(req.LotteryId))+":participants", strconv.Itoa(int(req.UserId)), participationJson).Err()
	if err != nil {
		logx.Error("参与抽奖失败")
		return nil, err
	}
	// 返回响应
	return
}
