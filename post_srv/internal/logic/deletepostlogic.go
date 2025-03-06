package logic

import (
	"LearningGuide/post_srv/internal/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "LearningGuide/post_srv/.PostProto"
	"LearningGuide/post_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePostLogic {
	return &DeletePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePostLogic) DeletePost(req *proto.DeletePostRequest) (*proto.Empty, error) {
	result := l.svcCtx.DB.Model(&model.Post{}).Delete(&model.Post{BaseModel: model.BaseModel{
		ID: req.Id,
	}})
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}

	return &proto.Empty{}, nil
}
