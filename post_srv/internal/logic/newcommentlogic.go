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

type NewCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNewCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewCommentLogic {
	return &NewCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NewCommentLogic) NewComment(req *proto.NewCommentRequest) (*proto.NewCommentResponse, error) {
	var post model.Post
	result := l.svcCtx.DB.Model(&model.Post{}).Where(&model.Post{
		BaseModel: model.BaseModel{
			ID: req.PostId,
		},
	}).Find(&post)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}

	comment := model.Comment{
		UserId:          req.UserId,
		PostId:          req.PostId,
		ParentCommentId: 0,
		Content:         req.Content,
	}

	var parentOwnerId int32

	if req.ParentCommentId != 0 {
		var parent model.Comment

		result := l.svcCtx.DB.Model(&model.Comment{}).Where(&model.Comment{
			BaseModel: model.BaseModel{
				ID: req.ParentCommentId,
			},
		}).Find(&parent)
		if result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "无效的回复评论ID")
		} else if parent.PostId != req.PostId {
			return nil, status.Errorf(codes.InvalidArgument, "与原评论不在同一帖子")
		} else {
			comment.ParentCommentId = req.ParentCommentId
			parentOwnerId = parent.UserId
		}
	}

	tx := l.svcCtx.DB.Begin()

	result = tx.Create(&comment)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	result = tx.Model(&model.Post{}).
		Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).
		Update("comment_num", post.CommentNum+1)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	if err := tx.Create(&model.Notice{
		UserId:  req.UserId,
		PostId:  req.PostId,
		OwnerId: post.UserId,
		Type:    model.NoticeTypeCommentToPost,
		IsRead:  false,
	}).Error; err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	if parentOwnerId != 0 {
		if err := tx.Create(&model.Notice{
			UserId:    req.UserId,
			PostId:    req.PostId,
			OwnerId:   parentOwnerId,
			Type:      model.NoticeTypeCommentToComment,
			CommentId: comment.ParentCommentId,
			IsRead:    false,
		}).Error; err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, result.Error.Error())
		}
	}

	tx.Commit()

	return &proto.NewCommentResponse{Id: comment.ID}, nil
}
