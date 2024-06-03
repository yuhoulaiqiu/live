package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	KqConsumerConf kq.KqConf
	Mysql          struct {
		DataSource string
	}
	Telemetry1 struct {
		Name     string
		Endpoint string
		Sampler  float64
		Batcher  string
	}
}
