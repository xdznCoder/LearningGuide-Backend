package logic

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/retriever"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"LearningGuide/chat_srv/.ChatProto"
	"LearningGuide/chat_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendStreamMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendStreamMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendStreamMessageLogic {
	return &SendStreamMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendStreamMessageLogic) SendStreamMessage(in *__ChatProto.UserMessage, stream __ChatProto.Chat_SendStreamMessageServer) error {
	docs, err := l.svcCtx.RAG.Retriever.Retrieve(l.ctx, in.Content,
		retriever.WithIndex(fmt.Sprintf("%s-%d", l.svcCtx.Config.ChatModel.Index, in.CourseID)),
		retriever.WithTopK(l.svcCtx.Config.ChatModel.TopK),
	)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	fmt.Println(docs)
	return nil
}
