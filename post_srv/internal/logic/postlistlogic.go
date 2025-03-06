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

type PostListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPostListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostListLogic {
	return &PostListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PostListLogic) PostList(req *proto.PostFilterRequest) (*proto.PostListResponse, error) {
	var post []model.Post
	filter := l.svcCtx.DB.Model(&model.Post{})

	if req.UserId != 0 {
		filter = filter.Where(&model.Post{UserId: req.UserId})
	}

	if len(req.Category) != 0 {
		filter = filter.Where("category in (?)", req.Category)
	}

	if req.Title != "" {
		filter = filter.Where("title LIKE ?", "%"+req.Title+"%")
	}

	var count int64

	filter.Count(&count)

	result := filter.Scopes(model.Paginate(int(req.PageNum), int(req.PageSize))).Find(&post)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	var respList []*proto.PostItemResponse

	for _, v := range post {
		respList = append(respList, &proto.PostItemResponse{
			UserId:     v.UserId,
			Category:   v.Category,
			Title:      v.Title,
			Desc:       v.Desc,
			Image:      v.Image,
			Id:         v.ID,
			LikeNum:    v.LikeNum,
			FavNum:     v.FavNum,
			CommentNum: v.CommentNum,
		})
	}

	return &proto.PostListResponse{
		Total: int32(count),
		Data:  respList,
	}, nil
}
