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

type NewLikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNewLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewLikeLogic {
	return &NewLikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NewLikeLogic) NewLike(req *proto.NewLikeRequest) (*proto.Empty, error) {
	var post model.Post
	result := l.svcCtx.DB.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).Find(&post)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}

	var like model.Like
	result = l.svcCtx.DB.Model(&model.Like{}).Where(&model.Like{UserId: req.UserId, PostId: req.PostId}).Find(&like)
	if result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已点赞")
	}

	tx := l.svcCtx.DB.Begin()

	result = tx.Create(&model.Like{
		UserId: req.UserId,
		PostId: req.PostId,
	})
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	result = tx.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).Update("like_num", post.LikeNum+1)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	tx.Commit()
	return &proto.Empty{}, nil
}
