package logic

import (
	"LearningGuide/class_srv/internal/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "LearningGuide/class_srv/.ClassProto"
	"LearningGuide/class_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLessonListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLessonListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLessonListLogic {
	return &GetLessonListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLessonListLogic) GetLessonList(in *proto.LessonFilterRequest) (*proto.LessonListResponse, error) {
	var RespList []*proto.LessonInfoResponse

	filter := l.svcCtx.DB.Preload("Course").Model(&model.Lesson{})

	if in.UserId != 0 {
		filter = filter.Where(&model.Lesson{UserId: in.UserId})
	}

	if in.DayOfWeek != 0 {
		filter = filter.Where(&model.Lesson{DayOfWeek: int(in.DayOfWeek)})
	}

	if in.WeekNum != 0 {
		filter = filter.Where(&model.Lesson{WeekNum: int(in.WeekNum)})
	}

	if in.CourseId != 0 {
		filter = filter.Where(&model.Lesson{CourseId: in.CourseId})
	}

	var lessons []model.Lesson

	result := filter.Find(&lessons)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	for _, v := range lessons {
		RespList = append(RespList, &proto.LessonInfoResponse{
			Id:        v.ID,
			UserId:    v.UserId,
			CourseId:  v.CourseId,
			WeekNum:   int32(v.WeekNum),
			DayOfWeek: int32(v.DayOfWeek),
			LessonNum: int32(v.LessonNum),
			Begin:     v.Begin,
			End:       v.End,
			Course: &proto.CourseInfoResponse{
				Id:       v.CourseId,
				Name:     v.Course.Name,
				Type:     v.Course.Type,
				CourseSn: v.Course.CourseSn,
				Image:    v.Course.Image,
				Teacher:  v.Course.Teacher,
				Credit:   v.Course.Credit,
			},
		})
	}

	return &proto.LessonListResponse{
		Total: int32(result.RowsAffected),
		Data:  RespList,
	}, nil
}
