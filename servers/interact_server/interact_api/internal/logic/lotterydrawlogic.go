package logic

import (
	"context"
	"live/models/interact_models"
	"live/servers/interact_server/interact_api/internal/svc"
	"live/servers/interact_server/interact_api/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type LotteryDrawLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLotteryDrawLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LotteryDrawLogic {
	return &LotteryDrawLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LotteryDrawLogic) LotteryDraw(req *types.LotteryRequest) (resp *types.LotteryResponse, err error) {
	// 将持续时间和开始时间解析为时间类型
	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		logx.Error("时间格式错误")
		return nil, err
	}

	loc, _ := time.LoadLocation("Local")
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, loc)
	if err != nil {
		logx.Error("时间格式错误")
		return nil, err
	}
	// 创建抽奖信息
	lottery := &interact_models.LotteryModel{
		AnchorId:      req.AnchorId,
		Prize:         req.Prize,
		Count:         req.Count,
		LotteryMethod: req.LotteryMethod,
		Duration:      duration,
		StartTime:     startTime,
		IsCompleted:   false,
	}

	// 保存抽奖信息到数据库
	l.svcCtx.DB.Create(&lottery)

	// 计算抽奖结束时间
	endTime := startTime.Add(duration)

	// 创建一个新的协程来监控当前时间
	go func() {
		for {
			// 每秒检查一次当前时间
			time.Sleep(time.Second)
			// 如果当前时间已经达到抽奖结束时间
			if time.Now().After(endTime) {
				// 更新数据库中的IsCompleted字段为true
				l.svcCtx.DB.Model(&interact_models.LotteryModel{}).Where("id = ?", lottery.ID).Update("IsCompleted", true)
				break
			}
		}
	}()

	return
}
