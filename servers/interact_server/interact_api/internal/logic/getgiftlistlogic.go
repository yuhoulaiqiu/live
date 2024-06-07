package logic

import (
	"context"
	"encoding/json"
	"live/servers/interact_server/interact_api/internal/svc"
	"live/servers/interact_server/interact_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGiftListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGiftListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGiftListLogic {
	return &GetGiftListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGiftListLogic) GetGiftList(req *types.GetGiftListRequest) (resp *types.GetGiftListResponse, err error) {
	//查redis中的gifts
	res, err := l.svcCtx.Redis.SMembers("gifts").Result()
	if err != nil {
		return nil, err
	}
	var gifts []types.GiftItem
	for _, v := range res {
		var gift types.GiftItem
		err := json.Unmarshal([]byte(v), &gift)
		if err != nil {
			return nil, err
		}
		gifts = append(gifts, gift)
	}
	return &types.GetGiftListResponse{
		Gifts: gifts,
	}, nil
}
