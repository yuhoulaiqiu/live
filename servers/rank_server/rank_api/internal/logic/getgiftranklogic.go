package logic

import (
	"context"
	"strconv"

	"live/servers/rank_server/rank_api/internal/svc"
	"live/servers/rank_server/rank_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGiftRankLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGiftRankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGiftRankLogic {
	return &GetGiftRankLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGiftRankLogic) GetGiftRank(req *types.GetGiftRequest) (resp *types.GetGiftResponse, err error) {
	result, err := l.svcCtx.Redis.ZRevRangeWithScores("gift_ranking", 0, req.TopN-1).Result()
	if err != nil {
		return nil, err
	}
	var giftRank []types.GiftItem
	for _, z := range result {
		id, _ := strconv.Atoi(z.Member.(string))

		giftRank = append(giftRank, types.GiftItem{
			AnchorID: id,
			Count:    int(z.Score),
		})
	}
	return &types.GetGiftResponse{
		Ranks: giftRank,
	}, nil
}
