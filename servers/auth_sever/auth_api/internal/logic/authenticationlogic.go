package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"live/servers/auth_sever/auth_api/internal/svc"
	"live/servers/auth_sever/auth_api/internal/types"
	"live/utils"
	"live/utils/jwts"
)

type AuthenticationLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	kafkaReader *kafka.Reader
}

func NewAuthenticationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthenticationLogic {
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "auth-requests",
		GroupID: "auth-service",
	})

	return &AuthenticationLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		kafkaReader: kafkaReader,
	}
}

func (l *AuthenticationLogic) Authentication(req *types.AuthenticationRequest) (resp *types.AuthenticationResponse, err error) {
	msg, err := l.kafkaReader.ReadMessage(context.Background())
	if err != nil {
		logx.Error("从 Kafka 读取信息失败")
		return nil, err
	}

	token := string(msg.Value)
	fmt.Println(token, req.ValidPath)
	if utils.InListByRegex(l.svcCtx.Config.WhiteList, req.ValidPath) {
		logx.Infof("%s 在白名单中", req.ValidPath)
		return
	}

	if token == "" {
		logx.Infof("token为空")
		err = errors.New("认证失败")
		return
	}
	claims, err := jwts.ParseToken(token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		logx.Error("token解析失败")
		err = errors.New("认证失败")
		return
	}
	_, err = l.svcCtx.Redis.Get(fmt.Sprintf("logout_%s", token)).Result()
	if err == nil {
		logx.Error("token已经登出")
		err = errors.New("认证失败")
		return
	}
	err = nil
	return &types.AuthenticationResponse{
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}
