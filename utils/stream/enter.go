package stream

import (
	"fmt"
	"os/exec"
)

func GetChannelKey(roomNumber string) (channelKey string) {
	//通过请求 http://localhost:8090/control/get?room=roomNubmer 获取channelKey
	output, err := exec.Command("curl", "http://localhost:8090/control/get?room="+roomNumber).Output()
	if err != nil {
		fmt.Printf("err:%v", err)
		return ""
	}
	str := string(output)
	//返回str的22到倒数第二位
	return str[22 : len(str)-2]

}
