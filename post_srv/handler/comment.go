package handler

import (
	"LearningGuide/post_srv/global"
	"LearningGuide/post_srv/model"
	proto "LearningGuide/post_srv/proto/.PostProto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p PostServer) NewComment(ctx context.Context, req *proto.NewCommentRequest) (*proto.NewCommentResponse, error) {
	var post model.Post
	result := global.DB.Model(&model.Post{}).Where(&model.Post{
		BaseModel: model.BaseModel{
			ID: req.PostId,
		},
	}).Find(&post)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的帖子ID")
	}

	comment := model.Comment{
		UserId:          req.UserId,
		PostId:          req.PostId,
		ParentCommentId: 0,
		Content:         req.Content,
	}

	if req.ParentCommentId != 0 {
		var parent model.Comment

		result := global.DB.Model(&model.Comment{}).Where(&model.Comment{
			BaseModel: model.BaseModel{
				ID: req.ParentCommentId,
			},
		}).Find(&parent)
		if result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "无效的回复评论ID")
		} else if parent.PostId != req.PostId {
			return nil, status.Errorf(codes.InvalidArgument, "与原评论不在同一帖子")
		} else {
			comment.ParentCommentId = req.ParentCommentId
		}
	}

	tx := global.DB.Begin()

	result = tx.Create(&comment)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	result = tx.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.PostId}}).Update("comment_num", post.CommentNum+1)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, result.Error.Error())
	}

	tx.Commit()

	return &proto.NewCommentResponse{Id: comment.ID}, nil
}

func (p PostServer) CommentList(ctx context.Context, req *proto.CommentFilterRequest) (*proto.CommentListResponse, error) {
	filter := global.DB.Model(&model.Comment{})

	if req.UserId != 0 {
		filter = filter.Where(&model.Comment{UserId: req.UserId})
	}

	if req.PostId != 0 && req.ParentCommendId == 0 {
		filter = filter.Where(&model.Comment{PostId: req.PostId, ParentCommentId: 0})
	} else if req.PostId == 0 && req.ParentCommendId != 0 {
		filter = filter.Where(&model.Comment{ParentCommentId: req.ParentCommendId})
	} else {
		filter = filter.Where(&model.Comment{PostId: req.PostId, ParentCommentId: req.ParentCommendId})
	}

	var count int64
	var comments []model.Comment

	filter.Count(&count)

	result := filter.Scopes(Paginate(int(req.PageNum), int(req.PageSize))).Find(&comments)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	var respList []*proto.CommentInfoResponse

	for _, v := range comments {
		respList = append(respList, &proto.CommentInfoResponse{
			UserId:          v.UserId,
			PostId:          v.PostId,
			ParentCommentId: v.ParentCommentId,
			Content:         v.Content,
			Id:              v.ID,
		})
	}

	return &proto.CommentListResponse{
		Total: int32(count),
		Data:  respList,
	}, nil
}

func (p PostServer) DeleteComment(ctx context.Context, req *proto.DeleteCommentRequest) (*proto.Empty, error) {
	var comment model.Comment
	var post model.Post

	result := global.DB.Model(&model.Comment{}).Where(&model.Comment{BaseModel: model.BaseModel{ID: req.Id}}).Find(&comment)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的评论ID")
	}

	global.DB.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: req.Id}}).Find(&post)

	tx := global.DB.Begin()

	result = tx.Model(&model.Comment{}).Delete(&model.Comment{BaseModel: model.BaseModel{ID: req.Id}})
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	result = tx.Model(&model.Comment{}).Where(&model.Comment{ParentCommentId: req.Id}).Delete(&model.Comment{})
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	subComment := int32(result.RowsAffected)

	result = tx.Model(&model.Post{}).Where(&model.Post{BaseModel: model.BaseModel{ID: comment.PostId}}).Update("comment_num", post.CommentNum-1-subComment)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, result.Error.Error())
	}

	tx.Commit()
	return &proto.Empty{}, nil
}
