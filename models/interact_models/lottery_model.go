package interact_models

import (
	"live/commen/models"
	"time"
)

type LotteryModel struct {
	models.Model
	AnchorId      uint          `json:"anchorId"`              // 主播ID
	Prize         string        `gorm:"size:256" json:"prize"` // 奖品
	Count         uint          `json:"count"`                 // 奖品数量
	LotteryMethod int           `json:"lotteryMethod"`         // 抽奖方式：0: 点击参与抽奖, 1: 发送弹幕抽奖, 2: 送礼物抽奖
	Duration      time.Duration `json:"duration"`              // 持续时间
	StartTime     time.Time     `json:"startTime"`             // 开始时间
	IsCompleted   bool          `json:"isCompleted"`           // 是否已完成
}
