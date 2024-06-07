package svc

import (
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-queue/kq"
	"gorm.io/gorm"
	"live/core"
	"live/servers/auth_sever/auth_api/internal/config"
	"time"
)

type ServiceContext struct {
	Config         config.Config
	DB             *gorm.DB
	Redis          *redis.Client
	KqPusherClient *kq.Pusher
}

func (s ServiceContext) Deadline() (deadline time.Time, ok bool) {
	//TODO implement me
	panic("implement me")
}

func (s ServiceContext) Done() <-chan struct{} {
	//TODO implement me
	panic("implement me")
}

func (s ServiceContext) Err() error {
	//TODO implement me
	panic("implement me")
}

func (s ServiceContext) Value(key any) any {
	//TODO implement me
	panic("implement me")
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitMysql(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Pwd, c.Redis.DB)
	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
		Redis:  client,
	}
}
