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

type NewPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNewPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewPostLogic {
	return &NewPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NewPostLogic) NewPost(req *proto.NewPostRequest) (*proto.NewPostResponse, error) {
	post := model.Post{
		UserId:     req.UserId,
		Category:   req.Category,
		Title:      req.Title,
		Content:    req.Content,
		Desc:       req.Desc,
		Image:      req.Image,
		LikeNum:    0,
		FavNum:     0,
		CommentNum: 0,
	}

	result := l.svcCtx.DB.Model(&model.Post{}).Create(&post)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &proto.NewPostResponse{Id: post.ID}, nil
}
