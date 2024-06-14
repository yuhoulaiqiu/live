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

//在rtmp连接建立以后，服务端要与客户端通过3次交换报文完成握手。
//与其他握手协议不同，rtmp协议握手交换的数据报文是固定大小的，客户端向服务端发送的3个报文为c0、c1、c2，
//服务端向客户端发送的3个报文为s0、s1、s2。c0与s0的大小为1个字节，c1与s1的大小为1536个字节，c2与s2的大小为1536个字节。
//发送顺序
//建立连接后，客户端开始发送C0、C1块到服务器；
//服务器端收到C0或C1后发送S0和S1；
//当客户端收齐S0和S1之后，开始发送C2；
//当服务端收齐C0和C1后，开发发送S2；
//当客户端收到S2，服务端收到C2，握手完成。
//在实际工程应用中，一般是客户端将C0、C1块同时发出，服务器在收到C1块之后同时将S0、S1、S2发给客户端。
//客户端收到S1之后，发送C2给服务端，握手完成。

// Handshake 进行 RTMP 握手
func (s *Session) Handshake() error {
	//C0和S0数据包占用一个字节，表示RTMP版本号。
	//目前RTMP版本定义为3,0-2是早期的专利产品所使用的值,现已经废弃,4-31是预留值,32-255是禁用值。

	// 第一步：接收 C0 和 C1
	c0c1 := make([]byte, 1537)
	if err := utils.ReadFullWithTimeout(s.conn, c0c1, 5*time.Second); err != nil {
		return err
	}
	if c0c1[0] != 0x03 {
		return errors.New("不支持的 RTMP 版本")
	}
	c1 := c0c1[1:]
	// 解析 C1
	c1Timestamp := c1[:4]

	// 第二步：发送 S0 和 S1
	s0s1 := make([]byte, 1537)
	s0s1[0] = 0x03

	//C1和S1数据包占用1536个字节。包含4个字节的时间戳，4个字节的0和1528个字节的随机数。
	// 填充 S1
	s1Timestamp := make([]byte, 4)
	copy(s1Timestamp, c1Timestamp) // 使用相同的时间戳
	copy(s0s1[1:5], s1Timestamp)
	copy(s0s1[5:9], make([]byte, 4)) // 4 个字节的 0
	if err := utils.FillRandomBytes(s0s1[9:1537]); err != nil {
		return err
	}
	if err := utils.WriteWithTimeout(s.conn, s0s1, 5*time.Second); err != nil {
		return err
	}
	//C2和S2数据包占用1536个字节，包含4个字节的时间戳，4个字节的对端的时间戳（C2数据包为S1数据包的时间戳，S2为C1数据包的时间戳）。
	// 第三步：接收 C2
	c2 := make([]byte, 1536)
	if err := utils.ReadFullWithTimeout(s.conn, c2, 5*time.Second); err != nil {
		return err
	}

	// 验证 C2
	if !utils.ValidateC2(s0s1[1:5], c2[:4]) {
		return errors.New("c2 验证失败")
	}

	// 第四步：发送 S2
	s2 := make([]byte, 1536)
	copy(s2[:4], s1Timestamp)  // S2 的时间戳为 S1 的时间戳
	copy(s2[4:8], c1Timestamp) // S2 的对端时间戳为 C1 的时间戳
	if err := utils.WriteWithTimeout(s.conn, s2, 5*time.Second); err != nil {
		return err
	}

	log.Println("RTMP 握手完成")
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
		return errors.New("没有可供发布的流媒体")
	}
	s.currentStream.Publisher = s
	log.Printf("Publishing to stream %d", s.currentStream.ID)
	return nil
}

// handlePlay 处理 play 命令
func (s *Session) handlePlay() error {
	if s.currentStream == nil {
		return errors.New("无播放流")
	}
	s.currentStream.AddSubscriber(s)
	log.Printf("Playing stream %d", s.currentStream.ID)
	return nil
}
