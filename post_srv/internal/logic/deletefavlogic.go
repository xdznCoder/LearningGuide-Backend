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

type DeleteFavLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFavLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFavLogic {
	return &DeleteFavLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFavLogic) DeleteFav(req *proto.DeleteFavRequest) (*proto.Empty, error) {
	var post model.Post
	result := l.svcCtx.DB.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).Find(&post)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}

	tx := l.svcCtx.DB.Begin()

	result = tx.Model(&model.Fav{}).Where(&model.Fav{
		UserId: req.UserId,
		PostId: req.PostId,
	}).Delete(&model.Fav{})
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "该用户尚未收藏")
	}

	result = tx.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).Update("fav_num", post.FavNum-1)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, result.Error.Error())
	}

	tx.Commit()

	return &proto.Empty{}, nil
}
