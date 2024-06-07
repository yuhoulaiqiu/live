package logic

import (
	"context"
	"strconv"

	"live/servers/rank_server/rank_api/internal/svc"
	"live/servers/rank_server/rank_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAnchorFansRankLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAnchorFansRankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAnchorFansRankLogic {
	return &GetAnchorFansRankLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAnchorFansRankLogic) GetAnchorFansRank(req *types.GetAnchorFansRankRequest) (resp *types.GetAnchorFansRankResponse, err error) {
	result, err := l.svcCtx.Redis.ZRevRangeWithScores("fans_ranking", 0, req.TopN-1).Result()
	if err != nil {
		return nil, err
	}

	var fansRank []types.AnchorItem
	for _, z := range result {
		id, _ := strconv.Atoi(z.Member.(string))
		fansRank = append(fansRank, types.AnchorItem{
			AnchorID: id,
			Fans:     int(z.Score),
		})
	}

	return &types.GetAnchorFansRankResponse{
		Ranks: fansRank,
	}, nil

}
