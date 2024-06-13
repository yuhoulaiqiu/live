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
	MsgTypeSetChunkSize  = 1
	MsgTypeAbortMessage  = 2
	MsgTypeAck           = 3
	MsgTypeUserControl   = 4
	MsgTypeWindowAckSize = 5
	MsgTypeSetPeerBw     = 6
	MsgTypeAudio         = 8
	MsgTypeVideo         = 9
	MsgTypeDataAMF0      = 18
	MsgTypeDataAMF3      = 15
	MsgTypeCommandAMF0   = 20
	MsgTypeCommandAMF3   = 17
	MsgTypeAggregate     = 22
)

// ChunkHeader 表示 RTMP Chunk Header
type ChunkHeader struct {
	Format          byte
	ChunkStreamID   uint32
	Timestamp       uint32
	MessageLength   uint32
	MessageTypeID   byte
	MessageStreamID uint32
}

// Message 表示 RTMP 消息
type Message struct {
	Header  ChunkHeader
	Payload []byte
}

// Protocol 表示 RTMP 协议处理器
type Protocol struct {
	conn      net.Conn
	chunkSize int
}

// NewProtocol 创建新的 RTMP 协议处理器
func NewProtocol(conn net.Conn) *Protocol {
	return &Protocol{
		conn:      conn,
		chunkSize: 128, // 默认Chunk大小
	}
}

// ReadMessage 读取 RTMP 消息
func (p *Protocol) ReadMessage() (*Message, error) {
	header, err := p.readChunkHeader()
	if err != nil {
		return nil, err
	}

	payload := make([]byte, header.MessageLength)
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
	header := make([]byte, 12)
	if _, err := io.ReadFull(p.conn, header); err != nil {
		return nil, err
	}

	format := header[0] >> 6
	chunkStreamID := uint32(header[0] & 0x3F)
	timestamp := readUint24(header[1:4])
	messageLength := readUint24(header[4:7])
	messageTypeID := header[7]
	messageStreamID := binary.LittleEndian.Uint32(header[8:12])

	return &ChunkHeader{
		Format:          format,
		ChunkStreamID:   chunkStreamID,
		Timestamp:       timestamp,
		MessageLength:   messageLength,
		MessageTypeID:   messageTypeID,
		MessageStreamID: messageStreamID,
	}, nil
}
func readUint24(b []byte) uint32 {
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// writeChunkHeader 写入 Chunk Header
func (p *Protocol) writeChunkHeader(header *ChunkHeader) error {
	buf := make([]byte, 12)
	buf[0] = (header.Format << 6) | byte(header.ChunkStreamID)
	writeUint24(buf[1:4], header.Timestamp)
	writeUint24(buf[4:7], header.MessageLength)
	buf[7] = header.MessageTypeID
	binary.LittleEndian.PutUint32(buf[8:12], header.MessageStreamID)

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

// HandleCommand 处理 RTMP 命令消息
func (p *Protocol) HandleCommand(msg *Message) error {
	// 简单处理 connect 命令
	if msg.Header.MessageTypeID == MsgTypeCommandAMF0 || msg.Header.MessageTypeID == MsgTypeCommandAMF3 {
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

// sendConnectResponse 发送 connect 命令的响应
func (p *Protocol) sendConnectResponse() error {
	response := []byte{
		0x02, 0x00, 0x07, // AMF0 string marker and length
		'r', 'e', 's', 'u', 'l', 't', 0x00, // "result"
		0x03,             // AMF0 object marker
		0x00, 0x00, 0x09, // object end marker
	}
	header := ChunkHeader{
		Format:          0,
		ChunkStreamID:   3,
		Timestamp:       0,
		MessageLength:   uint32(len(response)),
		MessageTypeID:   MsgTypeCommandAMF0,
		MessageStreamID: 0,
	}
	msg := &Message{
		Header:  header,
		Payload: response,
	}
	return p.WriteMessage(msg)
}
