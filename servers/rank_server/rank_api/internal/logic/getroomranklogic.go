package logic

import (
	"context"
	"live/servers/rank_server/rank_api/internal/svc"
	"live/servers/rank_server/rank_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoomRankLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoomRankLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoomRankLogic {
	return &GetRoomRankLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoomRankLogic) GetRoomRank(req *types.GetRoomRankRequest) (resp *types.GetRoomRankResponse, err error) {
	result, err := l.svcCtx.Redis.ZRevRangeWithScores("room_ranking", 0, req.TopN-1).Result()
	if err != nil {
		return nil, err
	}

	var roomRank []types.RoomItem
	for _, z := range result {
		roomRank = append(roomRank, types.RoomItem{
			RoomID:   z.Member.(string),
			Audience: int(z.Score),
		})
	}

	return &types.GetRoomRankResponse{
		Ranks: roomRank,
	}, nil
}
