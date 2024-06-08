package logic

import (
	"context"
	"live/models/interact_models"

	"live/servers/interact_server/interact_api/internal/svc"
	"live/servers/interact_server/interact_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatLogic) Chat(req *types.ChatRequest) (resp *types.ChatResponse, err error) {
	var chat interact_models.ChatModel
	chat.RoomNumber = req.RoomNumber
	chat.SendUserId = req.UserId
	chat.Msg = req.Content
	l.svcCtx.DB.Create(&chat)
	return
}
