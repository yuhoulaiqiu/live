package svc

import (
	"gorm.io/gorm"
	"live/core"
	"live/log_server/log_api/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Db     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDB := core.InitMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config: c,
		Db:     mysqlDB,
	}
}
