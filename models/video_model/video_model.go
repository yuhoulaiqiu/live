package video_model

import (
	"live/commen/models"
)

// VideoModel 视频/录播表
type VideoModel struct {
	models.Model
	Title      string `gorm:"size:64" json:"title"`      //视频标题
	Url        string `gorm:"size:256" json:"url"`       //视频地址
	Avatar     string `gorm:"size:256" json:"avatar"`    //封面
	PlayCount  int    `json:"playCount"`                 //播放次数
	Author     string `gorm:"size:32" json:"author"`     //作者
	FileFormat string `gorm:"size:32" json:"fileFormat"` //文件格式
}
