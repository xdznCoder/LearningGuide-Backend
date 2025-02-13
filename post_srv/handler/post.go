package handler

import (
	"LearningGuide/post_srv/global"
	"LearningGuide/post_srv/model"
	proto "LearningGuide/post_srv/proto/.PostProto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type PostServer struct {
	proto.UnimplementedPostServer
}

func (p PostServer) NewPost(ctx context.Context, req *proto.NewPostRequest) (*proto.NewPostResponse, error) {
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

	result := global.DB.Model(&model.Post{}).Create(&post)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &proto.NewPostResponse{Id: post.ID}, nil
}

func (p PostServer) GetPost(ctx context.Context, req *proto.PostID) (*proto.PostInfoResponse, error) {
	var post model.Post

	result := global.DB.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.Id}}).Find(&post)
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

func (p PostServer) PostList(ctx context.Context, req *proto.PostFilterRequest) (*proto.PostListResponse, error) {
	var post []model.Post
	filter := global.DB.Model(&model.Post{})

	if req.UserId != 0 {
		filter = filter.Where(&model.Post{UserId: req.UserId})
	}

	if req.Category != "" {
		filter = filter.Where(&model.Post{Category: req.Category})
	}

	if req.Title != "" {
		filter = filter.Where("title LIKE ?", "%"+req.Title+"%")
	}

	var count int64

	filter.Count(&count)

	result := filter.Scopes(Paginate(int(req.PageNum), int(req.PageSize))).Find(&post)
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

func (p PostServer) UpdatePost(ctx context.Context, req *proto.UpdatePostRequest) (*proto.Empty, error) {
	post := model.Post{
		Title:   req.Title,
		Content: req.Content,
		Desc:    req.Desc,
		Image:   req.Image,
	}

	result := global.DB.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.Id}}).Updates(&post)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}

	return &proto.Empty{}, nil
}

func (p PostServer) DeletePost(ctx context.Context, req *proto.DeletePostRequest) (*proto.Empty, error) {
	result := global.DB.Model(&model.Post{}).Delete(&model.Post{BaseModel: model.BaseModel{ID: req.Id}})
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}

	return &proto.Empty{}, nil
}

func Paginate(pageNum, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNum == 0 {
			pageNum = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
