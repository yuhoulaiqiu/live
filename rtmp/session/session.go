package session

import (
	"errors"
	"live/rtmp/protocol"
	"live/rtmp/stream"
	"live/rtmp/utils"
	"log"
	"net"
	"time"
)

// Session 表示一个 RTMP 会话
type Session struct {
	conn          net.Conn
	protocol      *protocol.Protocol
	streamManager *stream.StreamManager
	currentStream *stream.Stream
}

// NewSession 创建一个新的 RTMP 会话
func NewSession(conn net.Conn, streamManager *stream.StreamManager) *Session {
	return &Session{
		conn:          conn,
		protocol:      protocol.NewProtocol(conn),
		streamManager: streamManager,
	}
}

// GetConn 返回会话的连接
func (s *Session) GetConn() net.Conn {
	return s.conn
}

// Handshake 进行 RTMP 握手
func (s *Session) Handshake() error {
	// 第一步：接收 C0 和 C1
	c0c1 := make([]byte, 1537)
	if err := utils.ReadFullWithTimeout(s.conn, c0c1, 5*time.Second); err != nil {
		return err
	}
	if c0c1[0] != 0x03 {
		return errors.New("unsupported RTMP version")
	}
	c1 := c0c1[1:]

	// 第二步：发送 S0、S1 和 S2
	s0s1s2 := make([]byte, 3073)
	s0s1s2[0] = 0x03
	// 填充 S1
	if err := utils.FillRandomBytes(s0s1s2[1:1537]); err != nil {
		return err
	}
	// 填充 S2
	copy(s0s1s2[1537:], c1)
	if err := utils.WriteWithTimeout(s.conn, s0s1s2, 5*time.Second); err != nil {
		return err
	}

	// 第三步：接收 C2
	c2 := make([]byte, 1536)
	if err := utils.ReadFullWithTimeout(s.conn, c2, 5*time.Second); err != nil {
		return err
	}

	// 验证 C2 是否与 S1 匹配
	if !utils.ValidateC2(s0s1s2[1:1537], c2) {
		return errors.New("c2 validation failed")
	}

	log.Println("RTMP handshake completed")
	return nil
}

// HandleSession 处理 RTMP 会话
func (s *Session) HandleSession() {
	for {
		msg, err := s.protocol.ReadMessage()
		if err != nil {
			utils.LogError(err)
			return
		}

		if err := s.protocol.HandleCommand(msg); err != nil {
			utils.LogError(err)
			return
		}

		// 处理流命令
		if err := s.handleStreamCommand(msg); err != nil {
			utils.LogError(err)
			return
		}

		// 推送流数据给订阅者
		if s.currentStream != nil && s.currentStream.Publisher == s {
			s.currentStream.PushData(msg.Payload)
		}
	}
}

// handleStreamCommand 处理流命令
func (s *Session) handleStreamCommand(msg *protocol.Message) error {
	commandName, err := protocol.ParseCommandName(msg.Payload)
	if err != nil {
		return err
	}

	switch commandName {
	case "createStream":
		return s.handleCreateStream()
	case "publish":
		return s.handlePublish()
	case "play":
		return s.handlePlay()
	default:
		return nil
	}
}

// handleCreateStream 处理 createStream 命令
func (s *Session) handleCreateStream() error {
	streamID := uint32(1) //id先写死
	stream, err := s.streamManager.CreateStream(streamID)
	if err != nil {
		return err
	}
	s.currentStream = stream
	log.Printf("Stream %d created", streamID)
	return nil
}

// handlePublish 处理 publish 命令
func (s *Session) handlePublish() error {
	if s.currentStream == nil {
		return errors.New("no stream available for publishing")
	}
	s.currentStream.Publisher = s
	log.Printf("Publishing to stream %d", s.currentStream.ID)
	return nil
}

// handlePlay 处理 play 命令
func (s *Session) handlePlay() error {
	if s.currentStream == nil {
		return errors.New("no stream available for playing")
	}
	s.currentStream.AddSubscriber(s)
	log.Printf("Playing stream %d", s.currentStream.ID)
	return nil
}
