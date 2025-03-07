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

type NewFavLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNewFavLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewFavLogic {
	return &NewFavLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NewFavLogic) NewFav(req *proto.NewFavRequest) (*proto.Empty, error) {
	var post model.Post
	result := l.svcCtx.DB.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).Find(&post)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}

	var fav model.Fav
	result = l.svcCtx.DB.Model(&model.Fav{}).Where(&model.Fav{UserId: req.UserId, PostId: req.PostId}).Find(&fav)
	if result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已收藏")
	}

	tx := l.svcCtx.DB.Begin()

	// 创建数据库点赞信息
	result = tx.Create(&model.Fav{
		UserId: req.UserId,
		PostId: req.PostId,
	})
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	// 同步数据库帖子信息
	result = tx.Model(&model.Post{}).
		Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).
		Update("fav_num", post.FavNum+1)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	// 创建数据库通知信息
	if result := tx.Create(&model.Notice{
		UserId:  req.UserId,
		PostId:  req.PostId,
		OwnerId: post.UserId,
		Type:    model.NoticeTypeFavToPost,
		IsRead:  false,
	}); result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	tx.Commit()
	return &proto.Empty{}, nil
}
