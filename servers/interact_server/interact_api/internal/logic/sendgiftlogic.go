package logic

import (
	"context"
	"encoding/json"
	"errors"
	"live/models/user_models"
	"live/servers/interact_server/interact_api/internal/svc"
	"live/servers/interact_server/interact_api/internal/types"
	"strconv"
)

type SendGiftLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendGiftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendGiftLogic {
	return &SendGiftLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendGiftLogic) SendGift(req *types.SendGiftRequest) (*types.SendGiftResponse, error) {
	// 获取礼物信息
	res, err := l.svcCtx.Redis.SMembers("gifts").Result()
	if err != nil {
		return nil, err
	}
	var price float64
	for k, v := range res {
		if k == int(req.GiftID) {
			var gift types.GiftItem
			err := json.Unmarshal([]byte(v), &gift)
			if err != nil {
				return nil, err
			}
			price = float64(gift.Price)
			break
		}
	}

	// 开始数据库事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, err
	}

	// 获取用户余额
	var user user_models.UserModel
	tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", req.UserID).First(&user)

	// 检查用户余额是否足够
	if user.Balances < price*float64(req.Count) {
		tx.Rollback()
		return nil, errors.New("余额不足")
	}

	// 扣除用户余额
	user.Balances -= price * float64(req.Count)
	tx.Save(&user)

	// 获取主播信息
	var anchor user_models.UserModel
	tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", req.AnchorID).First(&anchor)

	// 增加主播余额
	anchor.Balances += price * 0.8 * float64(req.Count)
	tx.Save(&anchor)

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	// 将送礼记录写入redis
	//总记录：哪个直播间收益最高
	l.svcCtx.Redis.ZIncrBy("gift_ranking", price*float64(req.Count), strconv.Itoa(int(req.AnchorID)))
	//记录直播间内用户送礼物记录:榜一大哥
	l.svcCtx.Redis.ZIncrBy("gift_ranking_"+strconv.Itoa(int(req.AnchorID)), price*float64(req.Count), strconv.Itoa(int(req.UserID)))
	return &types.SendGiftResponse{}, nil

}
