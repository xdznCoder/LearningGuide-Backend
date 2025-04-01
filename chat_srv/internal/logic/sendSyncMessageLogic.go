package logic

import (
	"context"

	"LearningGuide/chat_srv/.ChatProto"
	"LearningGuide/chat_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendSyncMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendSyncMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSyncMessageLogic {
	return &SendSyncMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendSyncMessageLogic) SendSyncMessage(in *__ChatProto.UserMessage) (*__ChatProto.ChatModelResponse, error) {
	// todo: add your logic here and delete this line

	return &__ChatProto.ChatModelResponse{}, nil
}
