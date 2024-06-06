package live_models

import "live/common/models"

type GiftModel struct {
	models.Model
	Name   string `gorm:"size:64" json:"name"` //礼物名称
	Price  int    `json:"price"`               //价格
	Giver  int    `json:"giver"`               //赠送者
	Anchor int    `json:"anchor"`              //主播

}
