package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/sony/gobreaker"
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
	Addr string `json:"addr"`
	Etcd string `json:"etcd"`
	Log  logx.LogConf
}

var config Config

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *struct {
		UserID int `json:"userId"`
		Role   int `json:"role"`
	} `json:"data"`
}

// getService 从URL中获取服务名
func getService(url string) (string, error) {
	regex, _ := regexp.Compile(`/api/(.*?)/`)
	addrList := regex.FindStringSubmatch(url)
	if len(addrList) < 2 {
		return "", errors.New("invalid request")
	}
	return addrList[1], nil
}

// writeError 向响应中写入错误信息
func writeError(res http.ResponseWriter, msg string) {
	res.Write([]byte(msg))
}

// authenticate 发送认证请
func authenticate(req *http.Request, bodyBytes []byte) error {
	authAddr := etcd.GetServiceAddr(config.Etcd, "auth_api")
	authUrl := fmt.Sprintf("http://%s/api/auth/authentication", authAddr)
	authReq, err := http.NewRequest("POST", authUrl, bytes.NewReader(bodyBytes)) // 使用缓冲区的内容创建新的请求
	if err != nil {
		logx.Errorf("创建认证请求失败:%v", err)
		return err
	}
	authReq.Header = req.Header
	token := req.URL.Query().Get("token")
	if token != "" {
		authReq.Header.Set("Token", token)
	}
	authReq.Header.Set("ValidPath", req.URL.Path)
	authRes, err := http.DefaultClient.Do(authReq)
	if err != nil {
		logx.Errorf("认证请求失败:%v", err)
		return err
	}
	var authResponse Response
	byteData, _ := io.ReadAll(authRes.Body)
	err = json.Unmarshal(byteData, &authResponse)
	if err != nil {
		logx.Errorf("解析认证响应失败:%v", err)
		return err
	}
	fmt.Printf("认证响应:%v\n", authResponse)
	if authResponse.Code != 0 {
		return errors.New("验证失败")
	}

	// 设置请求头,返回认证数据
	if authResponse.Data != nil {
		req.Header.Set("User-ID", fmt.Sprintf("%d", authResponse.Data.UserID))
		req.Header.Set("Role", fmt.Sprintf("%d", authResponse.Data.Role))
	}
	return nil
}

type Proxy struct {
	limiter *rate.Limiter
	cb      *gobreaker.CircuitBreaker
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

	return &Proxy{
		limiter: rate.NewLimiter(r, b),
		cb:      cb,
	}
}

func (p *Proxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if !p.limiter.Allow() {
		http.Error(res, "请求太多", http.StatusTooManyRequests)
		return
	}

	//匹配请求前缀 /api/user/xxx
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
	// 使用断路器处理请求
	_, err = p.cb.Execute(func() (interface{}, error) {
		remoteAddr := strings.Split(req.RemoteAddr, ":")[0]
		logx.Infof("请求服务地址:%s, 客户端地址:%s", addr, remoteAddr)

		// 创建一个缓冲区并将req.Body的内容复制到其中
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 把原始的req.Body内容放回去

		// 请求认证服务地址
		if err := authenticate(req, bodyBytes); err != nil {
			writeError(res, error.Error(err))
			return nil, nil
		}

		// 创建一个新的请求并向目标服务发送请求
		remote, _ := url.Parse("http://" + addr)
		ReverseProxy := httputil.NewSingleHostReverseProxy(remote)
		ReverseProxy.ServeHTTP(res, req)
		return nil, nil
	})
	if err != nil {
		// 处理断路器返回的错误
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
	//绑定服务
	err := http.ListenAndServe(config.Addr, proxy)
	if err != nil {
		return
	}
}
