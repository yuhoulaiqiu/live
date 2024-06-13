package common

import "net"

// SessionInterface 定义了 Session 的接口
type SessionInterface interface {
	GetConn() net.Conn
}
