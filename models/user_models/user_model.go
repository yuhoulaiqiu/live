package user_models

import "live/common/models"

// UserModel 用户表
type UserModel struct {
	models.Model
	UserName string `gorm:"size:32" json:"userName"`
	NickName string `gorm:"size:32" json:"nickName"`
	Avatar   string `gorm:"size:256" json:"avatar"`
	Pwd      string `gorm:"size:64" json:"pwd"`
	Role     int8   `json:"role"` //1:管理员 2:普通用户
	Fans     int    `json:"fans"` //粉丝数
	//InWhich  string  `json:"inWhich"`                   //所在直播间
	Balances float64 `gorm:"default:0" json:"balances"` //余额
}
