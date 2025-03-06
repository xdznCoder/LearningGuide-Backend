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

type CreateLessonInBatchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLessonInBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLessonInBatchLogic {
	return &CreateLessonInBatchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLessonInBatchLogic) CreateLessonInBatch(in *proto.CreateLessonBatchRequest) (*proto.CreateLessonBatchResponse, error) {
	var ids []int32

	var lessons []model.Lesson

	if !beginEarlierThanEnd(in.Begin, in.End) {
		return nil, status.Errorf(codes.InvalidArgument, "无效的课时始末时间")
	}

	if in.EndWeek < in.BeginWeek {
		return nil, status.Errorf(codes.InvalidArgument, "无效的课程始末时间")
	}

	tx := l.svcCtx.DB.Begin()

	var course model.Course
	result := tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: in.CourseId}}).Find(&course)

	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "无效课程ID")
	}

	if course.UserId != in.UserId {
		tx.Rollback()
		return nil, status.Errorf(codes.InvalidArgument, "无效用户ID")
	}

	var former []model.Lesson

	tx.Model(&model.Lesson{}).Where(model.Lesson{
		UserId:    in.UserId,
		CourseId:  in.CourseId,
		Term:      course.Term,
		DayOfWeek: int(in.DayOfWeek),
		LessonNum: int(in.LessonNum),
	}).Find(&former)

	for _, v := range former {
		if int32(v.WeekNum) <= in.EndWeek && int32(v.WeekNum) >= in.BeginWeek {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "存在课时时间冲突")
		}
	}

	for i := in.BeginWeek; i <= in.EndWeek; i++ {
		lessons = append(lessons, model.Lesson{
			CourseId:  in.CourseId,
			UserId:    in.UserId,
			Term:      course.Term,
			WeekNum:   int(i),
			DayOfWeek: int(in.DayOfWeek),
			LessonNum: int(in.LessonNum),
			Begin:     in.Begin,
			End:       in.End,
		})
	}

	result = tx.Model(model.Lesson{}).Create(lessons)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.InvalidArgument, result.Error.Error())
	}

	result = tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: in.CourseId}}).
		Update("lesson_total", course.LessonTotal+in.EndWeek-in.BeginWeek+1)

	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	tx.Commit()

	for _, v := range lessons {
		ids = append(ids, v.ID)
	}

	return &proto.CreateLessonBatchResponse{Ids: ids}, nil
}
