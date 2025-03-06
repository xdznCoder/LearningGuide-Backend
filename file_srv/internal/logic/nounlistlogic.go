package logic

import (
	"LearningGuide/file_srv/internal/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "LearningGuide/file_srv/.FileProto"
	"LearningGuide/file_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type NounListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNounListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NounListLogic {
	return &NounListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NounListLogic) NounList(req *proto.NounListRequest) (*proto.NounListResponse, error) {
	if req.CourseId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "无效课程ID")
	}

	filter := l.svcCtx.DB.Model(&model.Noun{})

	if req.Name != "" {
		filter = filter.Where("name LIKE ?", "%"+req.Name+"%")
	}

	var count int64

	var nouns []model.Noun

	filter.Where(&model.Noun{CourseId: req.CourseId}).Count(&count)

	result := filter.Model(&model.Noun{}).Scopes(model.Paginate(int(req.PageNum), int(req.PageSize))).
		Where(&model.Noun{CourseId: req.CourseId}).Find(&nouns)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取名词列表出错: %v", result.Error)
	}

	var respList []*proto.NounInfoResponse

	for _, v := range nouns {
		respList = append(respList, &proto.NounInfoResponse{
			Id:       v.ID,
			Name:     v.Name,
			Content:  v.Content,
			CourseId: v.CourseId,
		})
	}

	return &proto.NounListResponse{
		Total: int32(count),
		Data:  respList,
	}, nil
}
