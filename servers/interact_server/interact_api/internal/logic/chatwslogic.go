package logic

import (
	"context"
	"live/servers/interact_server/interact_api/internal/svc"
	"live/servers/interact_server/interact_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatWsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatWsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatWsLogic) Chat(req *types.ChatRequest) (resp *types.ChatResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
