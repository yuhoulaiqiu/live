package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
	}
	Etcd  string
	Redis struct {
		Addr string
		Pwd  string
		DB   int
	}
	Telemetry1 struct {
		Name     string
		Endpoint string
		Sampler  float64
		Batcher  string
	}
}
