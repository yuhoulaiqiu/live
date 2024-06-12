package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Etcd      string
	FileSize  float64  // 文件大小限制
	WhiteList []string // 图片白名单
	UploadDir string   // 上传文件的目录
}
