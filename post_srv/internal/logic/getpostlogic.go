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

type GetPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLogic {
	return &GetPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPostLogic) GetPost(req *proto.PostID) (*proto.PostInfoResponse, error) {
	var post model.Post

	result := l.svcCtx.DB.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.Id}}).Find(&post)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}
	return &proto.PostInfoResponse{
		UserId:     post.UserId,
		Category:   post.Category,
		Content:    post.Content,
		Title:      post.Title,
		Desc:       post.Desc,
		Image:      post.Image,
		Id:         post.ID,
		LikeNum:    post.LikeNum,
		FavNum:     post.FavNum,
		CommentNum: post.CommentNum,
	}, nil
}
