package main

import (
	"fmt"
	"live/rtmp/session"
	"live/rtmp/stream"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	config := ServerConfig{
		Address: "0.0.0.0",
		Port:    1935,
	}

	server := NewServer(config)
	if err := server.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	// 运行一段时间后停止服务器
	time.Sleep(60 * time.Second)
	server.Stop()
}

// ServerConfig 包含服务器的配置
type ServerConfig struct {
	Address string // 监听地址
	Port    int    // 监听端口
}

// Server 表示 RTMP 服务器
type Server struct {
	config        ServerConfig
	listener      net.Listener
	wg            sync.WaitGroup
	quit          chan struct{}
	streamManager *stream.StreamManager
}

// NewServer 创建一个新的 RTMP 服务器
func NewServer(config ServerConfig) *Server {
	return &Server{
		config:        config,
		quit:          make(chan struct{}),
		streamManager: stream.NewStreamManager(),
	}
}

// Start 启动 RTMP 服务器
func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%d", s.config.Address, s.config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", address, err)
	}
	s.listener = listener
	log.Printf("RTMP server started on %s", address)

	s.wg.Add(1)
	go s.acceptConnections()

	return nil
}

// Stop 停止 RTMP 服务器
func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
	log.Println("RTMP server stopped")
}

// acceptConnections 接受客户端连接
func (s *Server) acceptConnections() {
	defer s.wg.Done()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Printf("failed to accept connection: %v", err)
				continue
			}
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

// handleConnection 处理客户端连接
func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	log.Printf("client connected: %s", conn.RemoteAddr().String())

	// 创建新的会话并进行握手
	session := session.NewSession(conn, s.streamManager)
	if err := session.Handshake(); err != nil {
		log.Printf("handshake failed: %v", err)
		return
	}

	// 处理 RTMP 会话
	session.HandleSession()
}

// ffmpeg -re -i demo.flv -c copy -f flv rtmp://localhost:1935/live/stream
