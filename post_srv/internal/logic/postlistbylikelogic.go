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

type PostListByLikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostListByLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostListByLikeLogic {
	return &PostListByLikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostListByLikeLogic) PostListByLike(req *proto.LikeListRequest) (*proto.PostListResponse, error) {
	var likes []model.Like
	var count int64

	l.svcCtx.DB.Model(&model.Like{}).Where(&model.Like{UserId: req.UserId}).Count(&count)
	result := l.svcCtx.DB.Model(&model.Like{}).Scopes(model.Paginate(int(req.PageNum), int(req.PageSize))).
		Where(&model.Like{UserId: req.UserId}).Find(&likes)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	var ids []int32

	for _, v := range likes {
		ids = append(ids, v.PostId)
	}

	var posts []model.Post

	result = l.svcCtx.DB.Model(model.Post{}).Where("id IN (?)", ids).Find(&posts)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	var respList []*proto.PostItemResponse

	for _, v := range posts {
		respList = append(respList, &proto.PostItemResponse{
			UserId:   v.UserId,
			Category: v.Category,
			Title:    v.Title,
			Desc:     v.Desc,
			Image:    v.Image,
			Id:       v.ID,
		})
	}

	return &proto.PostListResponse{
		Total: int32(count),
		Data:  respList,
	}, nil
}
