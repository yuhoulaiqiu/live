package mqs

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"live/servers/auth_sever/auth_api/internal/svc"
)

type PaymentSuccess struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPaymentSuccess(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentSuccess {
	return &PaymentSuccess{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentSuccess) Consume(key, val string) error {
	logx.Infof("PaymentSuccess key :%s , val :%s", key, val)
	return nil
}
