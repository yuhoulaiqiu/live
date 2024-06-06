package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/sony/gobreaker"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/time/rate"
	"io"
	"live/common/etcd"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var configFile = flag.String("f", "settings.yaml", "the config file")

type Config struct {
	Addr         string `json:"addr"`
	Etcd         string `json:"etcd"`
	Log          logx.LogConf
	KqPusherConf struct {
		Brokers []string
		Topic   string
	}
}

var config Config
var ctx = context.Background()

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *struct {
		UserID int `json:"userId"`
		Role   int `json:"role"`
	} `json:"data"`
}

func getService(url string) (string, error) {
	regex, _ := regexp.Compile(`/api/(.*?)/`)
	addrList := regex.FindStringSubmatch(url)
	if len(addrList) < 2 {
		return "", errors.New("invalid request")
	}
	return addrList[1], nil
}

func writeError(res http.ResponseWriter, msg string) {
	res.Write([]byte(msg))
}

func authenticate(req *http.Request, bodyBytes []byte, KqPusherClient *kq.Pusher) error {
	token := req.URL.Query().Get("token")
	err := KqPusherClient.Push(token)
	if err != nil {
		logx.Errorf("Failed to send message to Kafka: %v", err)
		return err
	}

	return nil
}

type Proxy struct {
	limiter        *rate.Limiter
	cb             *gobreaker.CircuitBreaker
	KqPusherClient *kq.Pusher
}

func NewProxy(r rate.Limit, b int) *Proxy {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "API Circuit Breaker",
		Timeout:     5 * time.Second,
		MaxRequests: 5,
		Interval:    60 * time.Second,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logx.Infof("断路器: %s 从 %s 到 %s\n", name, from, to)
		},
	})

	KqPusherClient := kq.NewPusher(config.KqPusherConf.Brokers, config.KqPusherConf.Topic)

	return &Proxy{
		limiter:        rate.NewLimiter(r, b),
		cb:             cb,
		KqPusherClient: KqPusherClient,
	}
}

func (p *Proxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if !p.limiter.Allow() {
		http.Error(res, "请求太多", http.StatusTooManyRequests)
		return
	}

	service, err := getService(req.URL.Path)
	if err != nil {
		writeError(res, "无效请求")
		return
	}

	addr := etcd.GetServiceAddr(config.Etcd, service+"_api")
	if addr == "" {
		logx.Errorf(" %s 不匹配的服务", service)
		writeError(res, "无效请求")
		return
	}

	_, err = p.cb.Execute(func() (interface{}, error) {
		remoteAddr := strings.Split(req.RemoteAddr, ":")[0]
		logx.Infof("请求服务地址:%s, 客户端地址:%s", addr, remoteAddr)

		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err := authenticate(req, bodyBytes, p.KqPusherClient); err != nil {
			writeError(res, err.Error())
			return nil, nil
		}

		remote, _ := url.Parse("http://" + addr)
		ReverseProxy := httputil.NewSingleHostReverseProxy(remote)
		ReverseProxy.ServeHTTP(res, req)
		return nil, nil
	})
	if err != nil {
		http.Error(res, "服务不可用", http.StatusServiceUnavailable)
		return
	}
}

func main() {
	flag.Parse()

	conf.MustLoad(*configFile, &config)
	logx.SetUp(config.Log)

	fmt.Printf("网关服务器运行在: %s\n", config.Addr)
	proxy := NewProxy(100, 50)
	err := http.ListenAndServe(config.Addr, proxy)
	if err != nil {
		return
	}
}
