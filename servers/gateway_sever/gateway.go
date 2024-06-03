package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"live/commen/etcd"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
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
	if authResponse.Code != 0 {
		return errors.New("authentication failed")
	}

	// 设置请求头,返回认证数据
	if authResponse.Data != nil {
		req.Header.Set("User-ID", fmt.Sprintf("%d", authResponse.Data.UserID))
		req.Header.Set("Role", fmt.Sprintf("%d", authResponse.Data.Role))
	}
	return nil
}

type Proxy struct {
}

func (Proxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	//匹配请求前缀 /api/user/xxx
	service, err := getService(req.URL.Path)
	if err != nil {
		writeError(res, "invalid request")
		return
	}

	addr := etcd.GetServiceAddr(config.Etcd, service+"_api")
	if addr == "" {
		logx.Errorf(" %s 不匹配的服务", service)
		writeError(res, "invalid request")
		return
	}

	remoteAddr := strings.Split(req.RemoteAddr, ":")[0]
	logx.Infof("请求服务地址:%s, 客户端地址:%s", addr, remoteAddr)

	// 创建一个缓冲区并将req.Body的内容复制到其中
	bodyBytes, _ := io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 把原始的req.Body内容放回去

	// 请求认证服务地址
	if err := authenticate(req, bodyBytes); err != nil {
		writeError(res, error.Error(err))
		return
	}

	// 创建一个新的请求并向目标服务发送请求
	remote, _ := url.Parse("http://" + addr)
	ReverseProxy := httputil.NewSingleHostReverseProxy(remote)
	ReverseProxy.ServeHTTP(res, req)
}

func main() {
	flag.Parse()

	conf.MustLoad(*configFile, &config)
	logx.SetUp(config.Log)

	fmt.Printf("gateway server running at %s\n", config.Addr)
	proxy := Proxy{}
	//绑定服务
	err := http.ListenAndServe(config.Addr, proxy)
	if err != nil {
		return
	}
}
