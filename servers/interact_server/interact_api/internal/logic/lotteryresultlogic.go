package logic

import (
	"context"
	"encoding/json"
	"errors"
	"live/models/interact_models"
	"math/rand"
	"strconv"
	"time"

	"live/servers/interact_server/interact_api/internal/svc"
	"live/servers/interact_server/interact_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LotteryResultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLotteryResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LotteryResultLogic {
	return &LotteryResultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LotteryResultLogic) LotteryResult(req *types.LotteryResultRequest) (resp *types.LotteryResultResponse, err error) {
	// 查找抽奖信息，确保抽奖已结束
	var lottery interact_models.LotteryModel
	if err := l.svcCtx.DB.Where("id = ? AND is_completed = ?", req.LotteryId, true).First(&lottery).Error; err != nil {
		logx.Error("抽奖不存在或尚未结束")
		return nil, errors.New("抽奖不存在或尚未结束")
	}
	// 从 Redis 中获取所有参与者
	participants, err := l.svcCtx.Redis.HGetAll("lottery:" + strconv.Itoa(int(req.LotteryId)) + ":participants").Result()
	if err != nil {
		logx.Error("查询参与用户失败")
		return nil, err
	}
	// 将参与者信息从 JSON 转换为结构体
	var participationModels []interact_models.LotteryParticipationModel
	for _, participantJson := range participants {
		var participation interact_models.LotteryParticipationModel
		err = json.Unmarshal([]byte(participantJson), &participation)
		if err != nil {
			logx.Error("查询参与用户失败")
			return nil, err
		}
		participationModels = append(participationModels, participation)
	}
	// 检查参与人数是否足够, 如果不足则返回参与的人
	if len(participants) < int(lottery.Count) {
		var users []types.Winner
		for _, participant := range participationModels {
			users = append(users, types.Winner{
				UserId: participant.UserId,
				Prize:  lottery.Prize,
			})
		}
		return &types.LotteryResultResponse{
			Winners: users,
		}, nil
	}

	// 随机挑选获奖者
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(participants), func(i, j int) {
		participationModels[i], participationModels[j] = participationModels[j], participationModels[i]
	})
	winners := participationModels[:lottery.Count]

	// 存储中奖结果
	var resultModels []interact_models.LotteryResultModel
	for _, winner := range winners {
		result := interact_models.LotteryResultModel{
			LotteryId: lottery.ID,
			UserId:    winner.UserId,
			Prize:     lottery.Prize,
		}
		resultModels = append(resultModels, result)
	}
	if err := l.svcCtx.DB.Create(&resultModels).Error; err != nil {
		logx.Error("存储中奖结果失败")
		return nil, err
	}
	// 从 Redis 中删除所有参与者
	err = l.svcCtx.Redis.Del("lottery:" + strconv.Itoa(int(req.LotteryId)) + ":participants").Err()
	if err != nil {
		logx.Error("删除参与者信息失败")
		return nil, err
	}
	// 构建响应
	var winnerResponses []types.Winner
	for _, result := range resultModels {
		winnerResponses = append(winnerResponses, types.Winner{
			UserId: result.UserId,
			Prize:  result.Prize,
		})
	}

	resp = &types.LotteryResultResponse{
		Winners: winnerResponses,
	}

	return resp, nil
}
