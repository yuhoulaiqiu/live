package main

import (
	"flag"
	"fmt"
	"live/common/etcd"
	"live/common/middleware"
	"live/servers/auth_sever/auth_api/internal/config"
	"live/servers/auth_sever/auth_api/internal/handler"
	"live/servers/auth_sever/auth_api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/auth.yaml", "the config file")

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
