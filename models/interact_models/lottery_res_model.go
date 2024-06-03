package interact_models

import "live/commen/models"

type LotteryResultModel struct {
	models.Model
	LotteryId uint   `json:"lotteryId"`             // 抽奖ID
	UserId    uint   `json:"userId"`                // 中奖用户ID
	Prize     string `gorm:"size:256" json:"prize"` // 奖品
}
