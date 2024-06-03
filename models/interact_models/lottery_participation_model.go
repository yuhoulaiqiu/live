package interact_models

import "live/commen/models"

type LotteryParticipationModel struct {
	models.Model
	LotteryId  uint `json:"lotteryId"`  // 抽奖ID
	UserId     uint `json:"userId"`     // 用户ID
	MethodType int  `json:"methodType"` // 参与方式：0: 点击参与抽奖, 1: 发送弹幕抽奖, 2: 送礼物抽奖
}
