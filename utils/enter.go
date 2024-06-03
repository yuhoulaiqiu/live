package utils

import (
	"crypto/md5"
	"github.com/zeromicro/go-zero/core/logx"
	"regexp"
)

func InList(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

// InListByRegex 正则匹配
func InListByRegex(list []string, item string) bool {
	for _, i := range list {
		regex, err := regexp.Compile(i)
		if err != nil {
			logx.Error(err)
			return false
		}
		if regex.MatchString(item) {
			return true
		}
	}
	return false
}

func MD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	cipherStr := h.Sum(nil)
	return string(cipherStr)
}
