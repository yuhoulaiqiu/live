package logic

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"live/servers/live_sever/live_api/internal/svc"
)

var CRooms = make(map[string]map[*websocket.Conn]bool)
var CLock = sync.RWMutex{}

type ConferenceWsLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	broadcast map[string]chan interface{}
	pc        map[string]*webrtc.PeerConnection // 每个房间一个PeerConnection
}

func NewConferenceWsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConferenceWsLogic {
	return &ConferenceWsLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		broadcast: make(map[string]chan interface{}),
		pc:        make(map[string]*webrtc.PeerConnection),
	}
}

func (l *ConferenceWsLogic) HandleConnections(wr http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(wr, r, nil)
	if err != nil {
		l.Logger.Error("Failed to upgrade to websocket:", err)
		return
	}
	defer conn.Close()

	roomNumber := r.URL.Query().Get("room")
	if roomNumber == "" {
		l.Logger.Error("未提供房间号")
		return
	}

	CLock.Lock()
	if _, ok := CRooms[roomNumber]; !ok {
		CRooms[roomNumber] = make(map[*websocket.Conn]bool)
	}
	CRooms[roomNumber][conn] = true
	CLock.Unlock()

	go l.readMessages(conn, roomNumber)
	l.writeMessages(conn, roomNumber)
}

func (l *ConferenceWsLogic) readMessages(conn *websocket.Conn, roomNumber string) {
	defer func() {
		CLock.Lock()
		delete(CRooms[roomNumber], conn)
		CLock.Unlock()
		conn.Close()
	}()

	for {
		var message map[string]interface{}
		err := conn.ReadJSON(&message)
		if err != nil {
			l.Logger.Error("Error reading json:", err)
			break
		}

		action, ok := message["action"].(string)
		if !ok {
			l.Logger.Error("Invalid action")
			continue
		}

		switch action {
		case "offer":
			l.handleOffer(conn, roomNumber, message)
		case "answer":
			l.handleAnswer(conn, roomNumber, message)
		case "candidate":
			l.handleCandidate(conn, roomNumber, message)
		}
	}
}

func (l *ConferenceWsLogic) writeMessages(conn *websocket.Conn, roomNumber string) {
	for {
		select {
		case msg := <-l.broadcast[roomNumber]:
			err := conn.WriteJSON(msg)
			if err != nil {
				l.Logger.Error("Error writing json:", err)
				conn.Close()
				return
			}
		}
	}
}

func (l *ConferenceWsLogic) handleOffer(conn *websocket.Conn, roomNumber string, message map[string]interface{}) {
	offer := webrtc.SessionDescription{}
	err := json.Unmarshal([]byte(message["offer"].(string)), &offer)
	if err != nil {
		l.Logger.Error("Error unmarshalling offer:", err)
		return
	}

	// 创建一个新的PeerConnection
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		// 如果创建PeerConnection失败，记录错误并返回
		l.Logger.Error("Error creating peer connection:", err)
		return
	}

	// 将PeerConnection存储在map中，使用房间号作为键
	l.pc[roomNumber] = pc

	// 当发现新的ICE候选时设置处理程序
	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c != nil {
			candidate, err := json.Marshal(c.ToJSON())
			if err != nil {
				l.Logger.Error("Error marshalling candidate:", err)
				return
			}
			// 创建一个包含候选人的消息
			message := map[string]interface{}{
				"action":    "candidate",
				"candidate": string(candidate),
			}
			// 将消息发送到房间的广播通道
			l.broadcast[roomNumber] <- message
		}
	})

	// 当收到新的媒体轨道时的处理函数
	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		// 记录收到的媒体轨道类型
		l.Logger.Infof("Got track: %s", track.Kind().String())

		// 创建一个本地的静态RTP媒体轨道，其参数与收到的远程媒体轨道相同
		localTrack, err := webrtc.NewTrackLocalStaticRTP(track.Codec().RTPCodecCapability, track.ID(), track.StreamID())
		if err != nil {
			// 如果创建本地媒体轨道失败，记录错误并返回
			l.Logger.Error("Error creating local track:", err)
			return
		}

		// 读取房间信息，准备向房间内的其他连接发送媒体轨道
		CLock.RLock()
		for wsConn := range CRooms[roomNumber] {
			// 不向发送offer的连接发送媒体轨道
			if wsConn != conn {
				pc, ok := l.pc[roomNumber]
				if ok {
					// 将本地媒体轨道添加到PeerConnection
					rtpSender, err := pc.AddTrack(localTrack)
					if err != nil {
						l.Logger.Error("Error adding track to peer connection:", err)
						continue
					}

					// 在新的goroutine中读取RTP数据并发送
					go func() {
						buf := make([]byte, 1500)
						for {
							if _, _, rtcpErr := rtpSender.Read(buf); rtcpErr != nil {
								return
							}
						}
					}()
				}
			}
		}
		CLock.RUnlock()

		// 在新的goroutine中读取远程媒体轨道的RTP数据
		go func() {
			buf := make([]byte, 1500)
			for {
				if _, _, rtcpErr := track.Read(buf); rtcpErr != nil {
					return
				}
			}
		}()
	})

	// 设置远程描述，这里的offer是从对方那里接收到的
	err = pc.SetRemoteDescription(offer)
	if err != nil {
		// 如果设置失败，记录错误并返回
		l.Logger.Error("Error setting remote description:", err)
		return
	}

	// 创建应答
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		// 如果创建应答失败，记录错误并返回
		l.Logger.Error("Error creating answer:", err)
		return
	}

	// 设置本地描述，这里的answer是我们自己创建的
	err = pc.SetLocalDescription(answer)
	if err != nil {
		// 如果设置失败，记录错误并返回
		l.Logger.Error("Error setting local description:", err)
		return
	}

	answerJSON, err := json.Marshal(answer)
	if err != nil {
		l.Logger.Error("Error marshalling answer:", err)
		return
	}

	// 创建一个包含应答的消息
	message = map[string]interface{}{
		"action": "answer",
		"answer": string(answerJSON),
	}
	// 将消息发送到房间的广播通道
	l.broadcast[roomNumber] <- message
}

// handleAnswer 处理从客户端接收到的应答
func (l *ConferenceWsLogic) handleAnswer(conn *websocket.Conn, roomNumber string, message map[string]interface{}) {
	answer := webrtc.SessionDescription{}
	err := json.Unmarshal([]byte(message["answer"].(string)), &answer)
	if err != nil {
		l.Logger.Error("Error unmarshalling answer:", err)
		return
	}

	// 从map中获取对应房间号的PeerConnection
	pc, ok := l.pc[roomNumber]
	if !ok {
		l.Logger.Error("PeerConnection not found for room:", roomNumber)
		return
	}

	// 设置远程描述，这里的answer是从对方那里接收到的
	err = pc.SetRemoteDescription(answer)
	if err != nil {
		l.Logger.Error("Error setting remote description:", err)
		return
	}
}

// handleCandidate 处理从客户端接收到的ICE候选信息
func (l *ConferenceWsLogic) handleCandidate(conn *websocket.Conn, roomNumber string, message map[string]interface{}) {
	candidate := webrtc.ICECandidateInit{}
	err := json.Unmarshal([]byte(message["candidate"].(string)), &candidate)
	if err != nil {
		l.Logger.Error("Error unmarshalling candidate:", err)
		return
	}

	// 从map中获取对应房间号的PeerConnection
	pc, ok := l.pc[roomNumber]
	if !ok {
		l.Logger.Error("PeerConnection not found for room:", roomNumber)
		return
	}

	// 将候选信息添加到PeerConnection
	err = pc.AddICECandidate(candidate)
	if err != nil {
		l.Logger.Error("Error adding ICE candidate:", err)
		return
	}
}
