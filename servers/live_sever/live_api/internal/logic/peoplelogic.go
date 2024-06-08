package logic

import (
	"context"
	"encoding/json"
	"errors"
	_ "fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"live/common/response"
	"live/models/live_models"
	"live/servers/live_sever/live_api/internal/svc"
	"net/http"
	"os/exec"
	"sync"
	"time"
)

// 直播间全局化，一个直播间对应一个ws连接
var rooms = make(map[string]map[*websocket.Conn]bool)
var Lock = sync.RWMutex{}

type WebRTCLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	broadcast map[string]chan interface{}
	pc        map[string]*webrtc.PeerConnection // 每个房间一个PeerConnection
}

func NewWebRTCLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebRTCLogic {
	return &WebRTCLogic{
		Logger:    logx.WithContext(ctx),
		ctx:       ctx,
		svcCtx:    svcCtx,
		broadcast: make(map[string]chan interface{}),
		pc:        make(map[string]*webrtc.PeerConnection),
	}
}

func (w *WebRTCLogic) HandleConnections(wr http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(wr, r, nil)
	if err != nil {
		logx.Error(err)
		response.Response(r, wr, nil, errors.New("websocket error"))
		return
	}

	defer conn.Close()

	roomNumber := r.URL.Query().Get("roomNumber")
	if roomNumber == "" {
		response.Response(r, wr, nil, errors.New("roomNumber is required"))
		return
	}

	// 判断房间号是否存在
	var liveModel live_models.LiveModel
	err = w.svcCtx.DB.Where("room_number = ?", roomNumber).First(&liveModel).Error
	if err != nil {
		response.Response(r, wr, nil, errors.New("房间号不存在"))
		return
	}

	Lock.Lock()
	if _, ok := rooms[roomNumber]; !ok {
		rooms[roomNumber] = make(map[*websocket.Conn]bool)
		w.broadcast[roomNumber] = make(chan interface{})
	}
	rooms[roomNumber][conn] = true
	Lock.Unlock()

	go w.readMessages(conn, roomNumber)
	go w.writeMessages(conn, roomNumber)
	go w.broadcastRankUpdates(roomNumber)

	select {}
}

func (w *WebRTCLogic) readMessages(conn *websocket.Conn, roomNumber string) {
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			// 用户断开连接，将其从客户端列表中移除
			Lock.Lock()
			delete(rooms[roomNumber], conn)
			Lock.Unlock()
			break
		}
		// 处理接收到的消息
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			logx.Error("Failed to unmarshal message:", err)
			continue
		}

		switch msg["type"] {
		case "offer":
			w.handleOffer(conn, roomNumber, msg)
		case "candidate":
			w.handleCandidate(conn, roomNumber, msg)
		}
	}
}

func (w *WebRTCLogic) writeMessages(conn *websocket.Conn, roomNumber string) {
	for {
		message := <-w.broadcast[roomNumber]
		for client := range rooms[roomNumber] {
			err := client.WriteJSON(message)
			if err != nil {
				logx.Error("websocket error:", err)
				client.Close()
				Lock.Lock()
				delete(rooms[roomNumber], client)
				Lock.Unlock()
			}
		}
	}
}

func (w *WebRTCLogic) broadcastRankUpdates(roomNumber string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 更新 Redis 中的实时人数
			w.svcCtx.Redis.ZAdd("room_ranking", redis.Z{
				Score:  float64(len(rooms[roomNumber])),
				Member: roomNumber,
			})
			// 发送实时人数
			w.broadcast[roomNumber] <- len(rooms[roomNumber])
		}
	}
}

func (w *WebRTCLogic) handleOffer(conn *websocket.Conn, roomNumber string, msg map[string]interface{}) {
	// 处理WebRTC Offer
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  msg["sdp"].(string),
	}
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		logx.Error("Failed to create PeerConnection:", err)
		return
	}

	w.pc[roomNumber] = pc

	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}
		resp := map[string]interface{}{
			"type":      "candidate",
			"candidate": candidate.ToJSON(),
		}
		w.broadcast[roomNumber] <- resp
	})

	err = pc.SetRemoteDescription(offer)
	if err != nil {
		logx.Error("Failed to set remote description:", err)
		return
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		logx.Error("Failed to create answer:", err)
		return
	}

	err = pc.SetLocalDescription(answer)
	if err != nil {
		logx.Error("Failed to set local description:", err)
		return
	}

	r := map[string]interface{}{
		"type": "answer",
		"sdp":  answer.SDP,
	}

	w.broadcast[roomNumber] <- r

	// 启动FFmpeg推流
	rtmpEndpoint := msg["rtmpEndpoint"].(string)
	go w.startFFmpegPush(roomNumber, rtmpEndpoint)
}

func (w *WebRTCLogic) handleCandidate(conn *websocket.Conn, roomNumber string, msg map[string]interface{}) {
	// 处理ICE候选
	SDPMid := msg["sdpMid"].(string)
	SDPMLineIndex := uint16(msg["sdpMLineIndex"].(float64))
	candidate := webrtc.ICECandidateInit{
		Candidate:     msg["candidate"].(string),
		SDPMid:        &SDPMid,
		SDPMLineIndex: &SDPMLineIndex,
	}

	pc, ok := w.pc[roomNumber]
	if !ok {
		logx.Error("PeerConnection not found for room:", roomNumber)
		return
	}

	err := pc.AddICECandidate(candidate)
	if err != nil {
		logx.Error("Failed to add ICE candidate:", err)
		return
	}
}

func (w *WebRTCLogic) startFFmpegPush(roomNumber, rtmpEndpoint string) {
	pc, ok := w.pc[roomNumber]
	if !ok {
		logx.Error("PeerConnection not found for room:", roomNumber)
		return
	}

	// 获取WebRTC轨道
	videoTrack := pc.GetTransceivers()[0].Receiver().Track()
	audioTrack := pc.GetTransceivers()[1].Receiver().Track()

	// 启动FFmpeg进程
	cmd := exec.Command("ffmpeg",
		"-i", "pipe:0",
		"-c:v", "libx264",
		"-preset", "veryfast",
		"-max_muxing_queue_size", "1024",
		"-f", "flv", rtmpEndpoint)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		logx.Error("Failed to get stdin pipe for FFmpeg:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		logx.Error("Failed to start FFmpeg process:", err)
		return
	}

	// 将WebRTC轨道写入FFmpeg stdin
	go func() {
		for {
			rtpPacket, _, err := videoTrack.ReadRTP()
			if err != nil {
				logx.Error("Failed to read RTP packet from video track:", err)
				break
			}
			_, err = stdin.Write(rtpPacket.Raw)
			if err != nil {
				logx.Error("Failed to write RTP packet to FFmpeg stdin:", err)
				break
			}
		}
	}()

	go func() {
		for {
			rtpPacket, _, err := audioTrack.ReadRTP()
			if err != nil {
				logx.Error("Failed to read RTP packet from audio track:", err)
				break
			}
			_, err = stdin.Write(rtpPacket.Raw)
			if err != nil {
				logx.Error("Failed to write RTP packet to FFmpeg stdin:", err)
				break
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		logx.Error("FFmpeg process exited with error:", err)
	}
}
