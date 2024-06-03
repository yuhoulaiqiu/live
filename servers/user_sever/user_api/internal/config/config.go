package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr string
		Pwd  string
		DB   int
	}
	Etcd       string
	Telemetry1 struct {
		Name     string
		Endpoint string
		Sampler  float64
		Batcher  string
	}
}
