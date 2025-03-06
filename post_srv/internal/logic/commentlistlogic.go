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

type CommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentListLogic) CommentList(req *proto.CommentFilterRequest) (*proto.CommentListResponse, error) {
	filter := l.svcCtx.DB.Model(&model.Comment{})

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

	result := filter.Scopes(model.Paginate(int(req.PageNum), int(req.PageSize))).Find(&comments)
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
