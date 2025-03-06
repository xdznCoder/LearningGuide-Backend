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

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCommentLogic) DeleteComment(req *proto.DeleteCommentRequest) (*proto.Empty, error) {
	var comment model.Comment
	var post model.Post

	result := l.svcCtx.DB.Model(&model.Comment{}).Where(&model.Comment{BaseModel: model.BaseModel{ID: req.Id}}).Find(&comment)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的评论ID")
	}

	l.svcCtx.DB.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.Id}}).Find(&post)

	tx := l.svcCtx.DB.Begin()

	result = tx.Model(&model.Comment{}).Delete(&model.Comment{BaseModel: model.BaseModel{ID: req.Id}})
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	result = tx.Model(&model.Comment{}).Where(&model.Comment{ParentCommentId: req.Id}).Delete(&model.Comment{})
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	subComment := int32(result.RowsAffected)

	result = tx.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: comment.PostId}}).Update("comment_num", post.CommentNum-1-subComment)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, result.Error.Error())
	}

	tx.Commit()
	return &proto.Empty{}, nil
}
