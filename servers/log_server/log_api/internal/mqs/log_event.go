package mqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"live/models/log_models"
	"live/models/user_models"
	"live/servers/log_server/log_api/internal/svc"
)

type LogEvent struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPaymentSuccess(ctx context.Context, svcCtx *svc.ServiceContext) *LogEvent {
	return &LogEvent{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type Request struct {
	LogType int8   `json:"logType"`
	IP      string `json:"IP"`
	UserID  uint   `json:"userID"`
	Addr    string `json:"addr"`
	Level   string `json:"level"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Service string `json:"service"`
}

func (l *LogEvent) Consume(key, val string) error {
	var req Request
	err := json.Unmarshal([]byte(val), &req)
	if err != nil {
		logx.Errorf("json解析错误, err :%v", err)
		return nil
	}
	// 查ip地址
	var info = log_models.LogModel{
		LogType: req.LogType,
		IP:      req.IP,
		UserID:  req.UserID,
		Level:   req.Level,
		Addr:    req.Addr,
		Title:   req.Title,
		Content: req.Content,
		Service: req.Service,
	}
	var user user_models.UserModel
	err = l.svcCtx.Db.Take(&user, req.UserID).Error
	if err != nil {
		logx.Errorf("用户不存在, err :%v", err)
		return nil
	}
	info.NickName = user.NickName
	info.Avatar = user.Avatar
	err = l.svcCtx.Db.Create(&info).Error
	if err != nil {
		logx.Errorf("写入日志错误, err :%v", err)
		return nil
	}

	return nil
}
