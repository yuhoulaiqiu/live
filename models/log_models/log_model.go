package log_models

import "live/common/models"

type LogModel struct {
	models.Model
	LogType  int8   `json:"log_type" gorm:"column:log_type;type:tinyint(1);not null;default:0;comment:'日志类型'"` //2:操作日志 3：运行日志
	IP       string `json:"ip" gorm:"column:ip;type:varchar(20);not null;default:'';comment:'ip地址'"`
	Addr     string `json:"addr" gorm:"column:addr;type:varchar(255);not null;default:'';comment:'地址'"`
	UserID   uint   `json:"user_id" gorm:"column:user_id;type:int(11);not null;default:0;comment:'用户id'"`
	NickName string `json:"nick_name" gorm:"column:nick_name;type:varchar(50);not null;default:'';comment:'用户昵称'"`
	Avatar   string `json:"avatar" gorm:"column:avatar;type:varchar(255);not null;default:'';comment:'用户头像'"`
	Level    string `json:"level" gorm:"column:level;type:varchar(20);not null;default:'';comment:'日志级别'"`
	Title    string `json:"title" gorm:"column:title;type:varchar(255);not null;default:'';comment:'标题'"`
	Content  string `json:"content" gorm:"column:content;type:text;not null;comment:'内容'"`
	Service  string `json:"service" gorm:"column:service;type:varchar(50);not null;default:'';comment:'服务名'"`
}
