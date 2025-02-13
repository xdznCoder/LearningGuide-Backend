package handler

import (
	"LearningGuide/post_srv/global"
	"LearningGuide/post_srv/model"
	proto "LearningGuide/post_srv/proto/.PostProto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p PostServer) NewFav(ctx context.Context, req *proto.NewFavRequest) (*proto.Empty, error) {
	var post model.Post
	result := global.DB.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).Find(&post)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}

	var fav model.Fav
	result = global.DB.Model(&model.Fav{}).Where(&model.Fav{UserId: req.UserId, PostId: req.PostId}).Find(&fav)
	if result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已收藏")
	}

	tx := global.DB.Begin()

	result = tx.Create(&model.Fav{
		UserId: req.UserId,
		PostId: req.PostId,
	})
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	result = tx.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).Update("fav_num", post.FavNum+1)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	tx.Commit()
	return &proto.Empty{}, nil
}

func (p PostServer) PostListByFav(ctx context.Context, req *proto.FavListRequest) (*proto.PostListResponse, error) {
	var favs []model.Fav
	var count int64

	global.DB.Model(&model.Fav{}).Where(&model.Fav{UserId: req.UserId}).Count(&count)
	result := global.DB.Model(&model.Fav{}).Scopes(Paginate(int(req.PageNum), int(req.PageSize))).
		Where(&model.Fav{UserId: req.UserId}).Find(&favs)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	var ids []int32

	for _, v := range favs {
		ids = append(ids, v.PostId)
	}

	var posts []model.Post

	result = global.DB.Model(model.Post{}).Where("id IN (?)", ids).Find(&posts)
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

func (p PostServer) DeleteFav(ctx context.Context, req *proto.DeleteFavRequest) (*proto.Empty, error) {
	var post model.Post
	result := global.DB.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).Find(&post)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}

	tx := global.DB.Begin()

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
