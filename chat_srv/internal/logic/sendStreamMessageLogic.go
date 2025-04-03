package logic

import (
	"LearningGuide/chat_srv/.ChatProto"
	FileProto "LearningGuide/chat_srv/.FileProto"
	"LearningGuide/chat_srv/internal/svc"
	"LearningGuide/chat_srv/internal/tpl"
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/duke-git/lancet/v2/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"

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
	var docs []*schema.Document
	var files []*schema.Document
	query := ""
	err := error(nil)

	switch tpl.TemplateType(in.TemplateType) {
	case tpl.TemplateTypeUserQuery:
		docs, err = l.svcCtx.RAG.Retriever.Retrieve(l.ctx, in.Content,
			retriever.WithIndex(fmt.Sprintf("%s-%d", l.svcCtx.Config.ChatModel.Index, in.CourseID)),
			retriever.WithTopK(l.svcCtx.Config.ChatModel.TopK),
		)

		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		query = in.Content
	case tpl.TemplateTypeExerciseGenerate:
		// 随机询问，防止重复出题
		docs, err = l.svcCtx.RAG.Retriever.Retrieve(l.ctx, fmt.Sprintf("%.10f", random.RandFloat(-1.0, 1.0, 10)),
			retriever.WithIndex(fmt.Sprintf("%s-%d", l.svcCtx.Config.ChatModel.Index, in.CourseID)),
			retriever.WithTopK(10),
		)
		if err != nil {
			l.Logger.Errorf("Redis索引出错：%v", err)
			return status.Error(codes.Internal, err.Error())
		}

		query = "请生成指定JSON格式的内容, 不要有任何多余的文字"
	case tpl.TemplateTypeMindMapGenerate:
		docs, err = l.svcCtx.RAG.LoadURLFile(l.ctx, in.FileURL)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "无效的URL链接")
		}

		query = "请生成JSMind格式的指定JSON内容, 不要有任何多余的文字"
	case tpl.TemplateTypeFileDescribeGenerate:
		docs, err = l.svcCtx.RAG.Loader.Load(l.ctx, document.Source{
			URI: in.FileURL,
		})
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "无效的URL链接")
		}

		query = "请描述文件的内容, 200字以内"
	case tpl.TemplateTypeNounExplainGenerate:
		docs, err = l.svcCtx.RAG.Retriever.Retrieve(l.ctx, in.Content,
			retriever.WithIndex(fmt.Sprintf("%s-%d", l.svcCtx.Config.ChatModel.Index, in.CourseID)),
			retriever.WithTopK(l.svcCtx.Config.ChatModel.TopK),
		)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		query = "请根据文档内容阐述名词的意思, 200字以内, 名词为 " + fmt.Sprintf("'%s'", query)
	default:
		return status.Error(codes.InvalidArgument, "无效的TemplateType")
	}

	var history []*schema.Message

	if in.TemplateType == int32(tpl.TemplateTypeUserQuery) {
		if in.SessionID == 0 {
			return status.Error(codes.InvalidArgument, "无效的会话ID")
		}
		resp, iErr := l.svcCtx.FileClient.MessageList(l.ctx, &FileProto.MessageListRequest{
			SessionId: in.SessionID,
			PageSize:  10,
			PageNum:   0,
		})
		if iErr != nil {
			return status.Error(codes.Internal, err.Error())
		}

		for _, h := range resp.Data {
			switch h.Speaker {
			case "user":
				if h.Type == 2 {
					f, loadErr := l.svcCtx.RAG.LoadURLFile(l.ctx, h.Content)
					if loadErr != nil {
						return status.Error(codes.Internal, loadErr.Error())
					}
					files = append(files, f...)
				}
				history = append(history, schema.UserMessage(h.Content))
			case "assistant":
				history = append(history, schema.AssistantMessage(h.Content, nil))
			}
		}
	}

	mes, err := tpl.Map.Get(l.ctx, tpl.TemplateMessage{
		Type:      tpl.TemplateType(in.TemplateType),
		Documents: docs,
		History:   history,
		File:      files,
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	output, err := l.svcCtx.RAG.Chatter.Stream(l.ctx, mes)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	fmt.Println(query)

	for {
		o, iErr := output.Recv()
		if iErr == io.EOF {
			output.Close()
			return nil
		}
		if iErr != nil {
			return status.Error(codes.Internal, iErr.Error())
		}
		fmt.Println("data: ", o.Content)
		iErr = stream.Send(&__ChatProto.ChatModelResponse{
			Content: o.Content,
		})
		if iErr != nil {
			return status.Error(codes.Internal, iErr.Error())
		}
	}
}
