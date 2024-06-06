package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"live/common/etcd"
	"live/common/middleware"
	"live/servers/live_sever/live_api/internal/config"
	"live/servers/live_sever/live_api/internal/handler"
	"live/servers/live_sever/live_api/internal/svc"
)

var configFile = flag.String("f", "etc/live.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	//设置全局中间件
	server.Use(middleware.LogMiddleware)
	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	server.Start()
}
