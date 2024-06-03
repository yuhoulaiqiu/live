package ips

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net"
)

func GetIP() (addr string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		logx.Errorf("获取本地IP失败:%v", err)
		return
	}
	for _, inter := range interfaces {
		addrs, err := inter.Addrs()
		if err != nil {
			logx.Errorf("获取本地IP失败:%v", err)
			continue
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					//检查ip是否为172.开头
					if ipNet.IP.String()[:3] == "172" {
						return ipNet.IP.String()
					}
				}
			}
		}
	}
	return
}
