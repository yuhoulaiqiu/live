package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"live/commen/middleware"
	"live/servers/log_server/log_api/internal/config"
	"live/servers/log_server/log_api/internal/handler"
	"live/servers/log_server/log_api/internal/mqs"
	"live/servers/log_server/log_api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/log.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	server.Use(middleware.LogMiddleware)
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range mqs.Consumers(c, context.Background(), ctx) {
		serviceGroup.Add(mq)
	}
	serviceGroup.Start()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
