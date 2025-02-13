package handler

import (
	"LearningGuide/post_srv/global"
	"LearningGuide/post_srv/model"
	proto "LearningGuide/post_srv/proto/.PostProto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p PostServer) NewLike(ctx context.Context, req *proto.NewLikeRequest) (*proto.Empty, error) {
	var post model.Post
	result := global.DB.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).Find(&post)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}

	var like model.Like
	result = global.DB.Model(&model.Like{}).Where(&model.Like{UserId: req.UserId, PostId: req.PostId}).Find(&like)
	if result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已点赞")
	}

	tx := global.DB.Begin()

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

func (p PostServer) PostListByLike(ctx context.Context, req *proto.LikeListRequest) (*proto.PostListResponse, error) {
	var likes []model.Like
	var count int64

	global.DB.Model(&model.Like{}).Where(&model.Like{UserId: req.UserId}).Count(&count)
	result := global.DB.Model(&model.Like{}).Scopes(Paginate(int(req.PageNum), int(req.PageSize))).
		Where(&model.Like{UserId: req.UserId}).Find(&likes)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	var ids []int32

	for _, v := range likes {
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

func (p PostServer) DeleteLike(ctx context.Context, req *proto.DeleteLikeRequest) (*proto.Empty, error) {
	var post model.Post
	result := global.DB.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).Find(&post)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}

	tx := global.DB.Begin()

	result = tx.Model(&model.Like{}).Where(&model.Like{
		UserId: req.UserId,
		PostId: req.PostId,
	}).Delete(&model.Like{})
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "该用户尚未点赞")
	}

	result = tx.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).Update("like_num", post.LikeNum-1)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, result.Error.Error())
	}

	tx.Commit()

	return &proto.Empty{}, nil
}
