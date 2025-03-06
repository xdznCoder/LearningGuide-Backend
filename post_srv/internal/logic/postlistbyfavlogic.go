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

type PostListByFavLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostListByFavLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostListByFavLogic {
	return &PostListByFavLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostListByFavLogic) PostListByFav(req *proto.FavListRequest) (*proto.PostListResponse, error) {
	var favs []model.Fav
	var count int64

	l.svcCtx.DB.Model(&model.Fav{}).Where(&model.Fav{UserId: req.UserId}).Count(&count)
	result := l.svcCtx.DB.Model(&model.Fav{}).Scopes(model.Paginate(int(req.PageNum), int(req.PageSize))).
		Where(&model.Fav{UserId: req.UserId}).Find(&favs)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	var ids []int32

	for _, v := range favs {
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
