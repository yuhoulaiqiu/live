package interact_models

import "live/common/models"

type GiftModel struct {
	models.Model
	Name  string  `gorm:"size:32" json:"name"`  //礼物名称
	Price float64 `json:"price"`                //价格
	Icon  string  `gorm:"size:256" json:"icon"` //图标
}
