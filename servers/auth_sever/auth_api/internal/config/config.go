package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Etcd  string
	Mysql struct {
		DataSource string
	}
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	Redis struct {
		Addr string
		Pwd  string
		DB   int
	}
	WhiteList    []string //白名单
	KqPusherConf struct {
		Brokers []string
		Topic   string
	}
	Telemetry1 struct {
		Name     string
		Endpoint string
		Sampler  float64
		Batcher  string
	}
	KqConsumerConf kq.KqConf
}
