package svc

import (
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-queue/kq"
	"gorm.io/gorm"
	"live/core"
	"live/servers/auth_sever/auth_api/internal/config"
)

type ServiceContext struct {
	Config         config.Config
	DB             *gorm.DB
	Redis          *redis.Client
	KqPusherClient *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitMysql(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Pwd, c.Redis.DB)
	return &ServiceContext{
		Config:         c,
		DB:             mysqlDb,
		Redis:          client,
		KqPusherClient: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
	}
}
