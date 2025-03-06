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

type UpdatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePostLogic {
	return &UpdatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePostLogic) UpdatePost(req *proto.UpdatePostRequest) (*proto.Empty, error) {
	post := model.Post{
		Title:   req.Title,
		Content: req.Content,
		Desc:    req.Desc,
		Image:   req.Image,
	}

	result := l.svcCtx.DB.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.Id}}).Updates(&post)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}

	return &proto.Empty{}, nil
}
