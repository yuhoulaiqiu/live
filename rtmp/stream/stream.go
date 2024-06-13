package stream

import (
	"errors"
	"live/rtmp/common"
	"sync"
)

// Stream 表示一个 RTMP 流
type Stream struct {
	ID          uint32
	Publisher   common.SessionInterface
	Subscribers []common.SessionInterface
	mutex       sync.Mutex
}

// NewStream 创建一个新的 RTMP 流
func NewStream(id uint32) *Stream {
	return &Stream{
		ID:          id,
		Subscribers: make([]common.SessionInterface, 0),
	}
}

// AddSubscriber 添加订阅者
func (s *Stream) AddSubscriber(sub common.SessionInterface) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Subscribers = append(s.Subscribers, sub)
}

// RemoveSubscriber 移除订阅者
func (s *Stream) RemoveSubscriber(sub common.SessionInterface) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i, subscriber := range s.Subscribers {
		if subscriber == sub {
			s.Subscribers = append(s.Subscribers[:i], s.Subscribers[i+1:]...)
			break
		}
	}
}

// PushData 推送数据给所有订阅者
func (s *Stream) PushData(data []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, subscriber := range s.Subscribers {
		conn := subscriber.GetConn()
		conn.Write(data)
	}
}

// StreamManager 管理所有 RTMP 流
type StreamManager struct {
	streams map[uint32]*Stream
	mutex   sync.Mutex
}

// NewStreamManager 创建一个新的 StreamManager
func NewStreamManager() *StreamManager {
	return &StreamManager{
		streams: make(map[uint32]*Stream),
	}
}

// CreateStream 创建新的流
func (sm *StreamManager) CreateStream(id uint32) (*Stream, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	if _, exists := sm.streams[id]; exists {
		return nil, errors.New("stream already exists")
	}
	stream := NewStream(id)
	sm.streams[id] = stream
	return stream, nil
}

// GetStream 获取指定 ID 的流
func (sm *StreamManager) GetStream(id uint32) (*Stream, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	stream, exists := sm.streams[id]
	if !exists {
		return nil, errors.New("stream not found")
	}
	return stream, nil
}

// DeleteStream 删除指定 ID 的流
func (sm *StreamManager) DeleteStream(id uint32) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	delete(sm.streams, id)
}
