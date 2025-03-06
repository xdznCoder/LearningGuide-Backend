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

type ExerciseListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExerciseListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExerciseListLogic {
	return &ExerciseListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ExerciseListLogic) ExerciseList(in *proto.ExerciseListRequest) (*proto.ExerciseListResponse, error) {
	filter := l.svcCtx.DB.Model(&model.Exercise{})

	if in.CourseId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "无效的课程ID")
	}

	if in.Question != "" {
		filter = filter.Where("question LIKE ?", "%"+in.Question+"%")
	}

	if in.IsRight != "" {
		switch in.IsRight {
		case "true":
			filter = filter.Where(&model.Exercise{IsRight: "true"})
		case "false":
			filter = filter.Where(&model.Exercise{IsRight: "false"})
		default:
			return nil, status.Errorf(codes.InvalidArgument, "无效的is_right参数")
		}
	}

	var count int64

	var exercises []model.Exercise

	filter.Count(&count)

	result := filter.Scopes(model.Paginate(int(in.PageNum), int(in.PageSize))).Find(&exercises)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	var respList []*proto.ExerciseInfoResponse

	for _, v := range exercises {
		respList = append(respList, &proto.ExerciseInfoResponse{
			CourseId: v.CourseId,
			Question: v.Question,
			SectionA: v.SectionA,
			SectionB: v.SectionB,
			SectionC: v.SectionC,
			SectionD: v.SectionD,
			Answer:   v.Answer,
			Reason:   v.Reason,
			IsRight:  v.IsRight,
			Id:       v.ID,
		})
	}

	return &proto.ExerciseListResponse{
		Total: int32(count),
		Data:  respList,
	}, nil
}
