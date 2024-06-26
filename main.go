package main

import (
	"flag"
	"fmt"
	"live/core"
	"live/models/interact_models"
	"live/models/live_models"
	"live/models/log_models"
	"live/models/user_models"
	"live/models/video_model"
)

type Option struct {
	DB bool
}

func main() {
	var opt Option
	flag.BoolVar(&opt.DB, "db", false, "初始化数据库")
	flag.Parse()
	mysqlDataSource := "root:zxc3240858086@tcp(127.0.0.1:3306)/live_db?charset=utf8mb4&parseTime=True&loc=Local"
	if opt.DB {
		db := core.InitMysql(mysqlDataSource)
		err := db.AutoMigrate(&user_models.UserModel{},
			&user_models.FansModel{},
			&live_models.LiveModel{},
			&interact_models.GiftModel{},
			&video_model.VideoModel{},
			&log_models.LogModel{},
			&interact_models.LotteryModel{},
			&interact_models.LotteryResultModel{},
			&interact_models.ChatModel{},
		)
		if err != nil {
			fmt.Println("数据库初始化失败", err)
			return
		}
		// 添加索引
		if err := db.Exec("CREATE INDEX idx_fans_relation ON fans_models (anchor_id, user_id)").Error; err != nil {
			return
		}
		fmt.Println("数据库初始化成功")
	}
}
