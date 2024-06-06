package live_models

import "live/common/models"

type LiveModel struct {
	models.Model
	Title         string `gorm:"size:64" json:"title"`        //直播标题
	RoomNumber    string `gorm:"size:32" json:"roomNumber"`   //房间号
	Avatar        string `gorm:"size:256" json:"avatar"`      //封面
	Description   string `gorm:"size:256" json:"description"` //描述
	AnchorId      uint   `json:"userId"`                      //主播id
	AudienceCount int    `json:"audienceCount"`               //观众数
	IsStart       bool   `json:"isStart"`                     //是否开始
	RTMPAddress   string `json:"rtmpAddress"`                 //rtmp地址
}
