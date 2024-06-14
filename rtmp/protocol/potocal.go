package protocol

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

// RTMP 消息类型
const (
	MsgTypeSetChunkSize  = 1  // 设置 Chunk 大小
	MsgTypeAbortMessage  = 2  // 中断消息
	MsgTypeAck           = 3  // 确认消息
	MsgTypeUserControl   = 4  // 用户控制消息
	MsgTypeWindowAckSize = 5  // 窗口确认大小
	MsgTypeSetPeerBw     = 6  // 设置对等带宽
	MsgTypeAudio         = 8  // 音频消息
	MsgTypeVideo         = 9  // 视频消息
	MsgTypeDataAMF0      = 18 // 数据消息（AMF0 格式）
	MsgTypeDataAMF3      = 15 // 数据消息（AMF3 格式）
	MsgTypeCommandAMF0   = 20 // 命令消息（AMF0 格式）
	MsgTypeCommandAMF3   = 17 // 命令消息（AMF3 格式）
	MsgTypeAggregate     = 22 // 聚合消息
)

//12字节的Chunk Header： 2bits:format,6bits:chunk stream id,3bytes:timestamp,3bytes:body size,1byte:type id,4bytes:stream id
//8字节的Chunk Header： 2bits:format,6bits:chunk stream id,3bytes:timestamp delta,3bytes:body size,1byte:type id
//4字节的Chunk Header： 2bits:format,6bits:chunk stream id,3bytes:timestamp delta
//1字节的Chunk Header： 2bits:format,6bits:chunk stream id

// ChunkHeader 表示 RTMP Chunk Header
type ChunkHeader struct {
	Format        byte   // Format决定了 RTMP header 的长度为多少个字节, 0:12, 1:8, 2:4, 3:1
	ChunkStreamID uint32 // 分块流 ID
	Timestamp     uint32 // 时间戳
	BodySize      uint32 // RTMP Body 所包含数据包的大小
	TypeID        byte   // 消息类型 ID, 0x14 表示以 AMF0 编码
	StreamID      uint32 // 消息流 ID
}

// Message 表示 RTMP 消息
type Message struct {
	Header  ChunkHeader
	Payload []byte
}

// Protocol 表示 RTMP 协议处理器
type Protocol struct {
	conn      net.Conn // 一种面向流的网络连接
	chunkSize int      // Chunk 大小
}

// NewProtocol 创建新的 RTMP 处理器
func NewProtocol(conn net.Conn) *Protocol {
	return &Protocol{
		conn:      conn,
		chunkSize: 128, // 默认 Chunk 大小
	}
}

// ReadMessage 读取 RTMP 消息
func (p *Protocol) ReadMessage() (*Message, error) {
	header, err := p.readChunkHeader()
	if err != nil {
		return nil, err
	}

	payload := make([]byte, header.BodySize)
	if _, err := io.ReadFull(p.conn, payload); err != nil {
		return nil, err
	}

	return &Message{
		Header:  *header,
		Payload: payload,
	}, nil
}

// WriteMessage 写入 RTMP 消息
func (p *Protocol) WriteMessage(msg *Message) error {
	if err := p.writeChunkHeader(&msg.Header); err != nil {
		return err
	}
	if _, err := p.conn.Write(msg.Payload); err != nil {
		return err
	}
	return nil
}

// readChunkHeader 读取 Chunk Header
func (p *Protocol) readChunkHeader() (*ChunkHeader, error) {
	header := make([]byte, 1)
	if _, err := io.ReadFull(p.conn, header); err != nil {
		return nil, err
	}

	format := header[0] >> 6                  // 前 2 位表示 Format
	chunkStreamID := uint32(header[0] & 0x3F) // 后 6 位表示 Chunk Stream ID

	var timestamp, bodySize, streamID uint32
	var typeID byte

	switch format {
	case 0:
		header = make([]byte, 11)
		if _, err := io.ReadFull(p.conn, header); err != nil {
			return nil, err
		}
		timestamp = readUint24(header[0:3])
		bodySize = readUint24(header[3:6])
		typeID = header[6]
		streamID = binary.LittleEndian.Uint32(header[7:11])
	case 1:
		header = make([]byte, 7)
		if _, err := io.ReadFull(p.conn, header); err != nil {
			return nil, err
		}
		timestamp = readUint24(header[0:3])
		bodySize = readUint24(header[3:6])
		typeID = header[6]
	case 2:
		header = make([]byte, 3)
		if _, err := io.ReadFull(p.conn, header); err != nil {
			return nil, err
		}
		timestamp = readUint24(header)
	case 3:
		// 格式 3 无需额外字节
	default:
		return nil, errors.New("invalid chunk header format")
	}

	return &ChunkHeader{
		Format:        format,
		ChunkStreamID: chunkStreamID,
		Timestamp:     timestamp,
		BodySize:      bodySize,
		TypeID:        typeID,
		StreamID:      streamID,
	}, nil
}

func readUint24(b []byte) uint32 {
	// 读取 3 字节的整数
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// writeChunkHeader 写入 Chunk Header
func (p *Protocol) writeChunkHeader(header *ChunkHeader) error {
	var buf []byte
	switch header.Format {
	case 0:
		buf = make([]byte, 12)
		buf[0] = (header.Format << 6) | byte(header.ChunkStreamID)
		writeUint24(buf[1:4], header.Timestamp)
		writeUint24(buf[4:7], header.BodySize)
		buf[7] = header.TypeID
		binary.LittleEndian.PutUint32(buf[8:12], header.StreamID)
	case 1:
		buf = make([]byte, 8)
		buf[0] = (header.Format << 6) | byte(header.ChunkStreamID)
		writeUint24(buf[1:4], header.Timestamp)
		writeUint24(buf[4:7], header.BodySize)
		buf[7] = header.TypeID
	case 2:
		buf = make([]byte, 4)
		buf[0] = (header.Format << 6) | byte(header.ChunkStreamID)
		writeUint24(buf[1:4], header.Timestamp)
	case 3:
		buf = make([]byte, 1)
		buf[0] = (header.Format << 6) | byte(header.ChunkStreamID)
	default:
		return errors.New("invalid chunk header format")
	}

	if _, err := p.conn.Write(buf); err != nil {
		return err
	}
	return nil
}

func writeUint24(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}

//当rtmp客户端和rtmp服务端握手完成之后，客户端就会向服务端发送connect消息。connect消息的格式按照RTMP Header+RTMP Body的格式组织。
//其中RTMP Header的Type ID为0x14，表示以AMF0编码的command消息。

// HandleCommand 处理 RTMP 命令消息
func (p *Protocol) HandleCommand(msg *Message) error {
	// 简单处理 connect 命令
	if msg.Header.TypeID == MsgTypeCommandAMF0 || msg.Header.TypeID == MsgTypeCommandAMF3 {
		commandName, err := ParseCommandName(msg.Payload)
		if err != nil {
			return err
		}
		if commandName == "connect" {
			log.Println("Received connect command")
			return p.sendConnectResponse()
		}
	}
	return nil
}

// ParseCommandName 解析命令名称
func ParseCommandName(payload []byte) (string, error) {
	// 简单解析 AMF0 命令名称
	if len(payload) < 3 {
		return "", errors.New("invalid command payload")
	}
	if payload[0] != 0x02 { // AMF0 string marker
		return "", errors.New("invalid command format")
	}
	length := binary.BigEndian.Uint16(payload[1:3])
	if len(payload) < int(3+length) {
		return "", errors.New("invalid command length")
	}
	return string(payload[3 : 3+length]), nil
}

// ParseCommandPayload 解析命令   待实现
func ParseCommandPayload(payload []byte) (uint32, any, error) {
	return 0, nil, nil
}

// ParsePublishPayload 解析发布命令的Payload   待实现
func ParsePublishPayload(payload []byte) (string, string, error) {
	return "", "", nil
}

func ParsePlayPayload(payload []byte) (string, int, int, bool, error) {
	return "", 0, 0, false, nil

}

// sendConnectResponse 发送 connect 命令的响应
func (p *Protocol) sendConnectResponse() error {
	response := []byte{
		0x02, 0x00, 0x07, // AMF0 字符串标记和长度
		'r', 'e', 's', 'u', 'l', 't', 0x00, // "result"
		0x03,             // AMF0 物体标记
		0x00, 0x00, 0x09, // 对象结束标记
	}
	header := ChunkHeader{
		Format:        0,
		ChunkStreamID: 3,
		Timestamp:     0,
		BodySize:      uint32(len(response)),
		TypeID:        MsgTypeCommandAMF0,
		StreamID:      0,
	}
	msg := &Message{
		Header:  header,
		Payload: response,
	}
	return p.WriteMessage(msg)
}
