package logic

import (
	"LearningGuide/chat_srv/.ChatProto"
	"LearningGuide/chat_srv/internal/svc"
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/document"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/url"
	"path/filepath"
)

type UploadDocumentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadDocumentLogic {
	return &UploadDocumentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadDocumentLogic) UploadDocument(in *__ChatProto.CourseDocument) (*__ChatProto.Empty, error) {
	parsedURL, err := url.Parse(in.URL)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "无效的URL链接")
	}
	fileName := filepath.Base(parsedURL.Path)

	err = l.svcCtx.RAG.InitVectorIndex(l.ctx, l.svcCtx.Redis, &l.svcCtx.Config, in.CourseID)
	if err != nil {
		l.Errorf("InitVectorIndex err:%v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	docs, err := l.svcCtx.RAG.Loader.Load(l.ctx, document.Source{
		URI: in.URL,
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "无效的URL链接")
	}

	docs, err = l.svcCtx.RAG.Splitter.Transform(l.ctx, docs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	for i, d := range docs {
		d.ID = fmt.Sprintf("Course-%d:%s-chunk%d", in.CourseID, fileName, i)
	}

	_, err = l.svcCtx.RAG.Indexer.Store(l.ctx, docs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &__ChatProto.Empty{}, nil
}
