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

type DeleteLikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLikeLogic {
	return &DeleteLikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLikeLogic) DeleteLike(req *proto.DeleteLikeRequest) (*proto.Empty, error) {
	var post model.Post
	result := l.svcCtx.DB.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).Find(&post)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}

	tx := l.svcCtx.DB.Begin()

	result = tx.Model(&model.Like{}).Where(&model.Like{
		UserId: req.UserId,
		PostId: req.PostId,
	}).Delete(&model.Like{})
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "该用户尚未点赞")
	}

	result = tx.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).Update("like_num", post.LikeNum-1)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, result.Error.Error())
	}

	tx.Commit()

	return &proto.Empty{}, nil
}
