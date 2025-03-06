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

type GetLessonDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLessonDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLessonDetailLogic {
	return &GetLessonDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLessonDetailLogic) GetLessonDetail(in *proto.LessonDetailRequest) (*proto.LessonInfoResponse, error) {

	var lesson model.Lesson

	result := l.svcCtx.DB.Preload("Course").Where(&model.Lesson{BaseModel: model.BaseModel{ID: in.Id}}).Find(&lesson)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效课时ID")
	}

	return &proto.LessonInfoResponse{
		Id:        lesson.ID,
		CourseId:  lesson.CourseId,
		WeekNum:   int32(lesson.WeekNum),
		DayOfWeek: int32(lesson.DayOfWeek),
		LessonNum: int32(lesson.LessonNum),
		Begin:     lesson.Begin,
		End:       lesson.End,
		Course: &proto.CourseInfoResponse{
			Id:          lesson.Course.ID,
			Name:        lesson.Course.Name,
			Type:        lesson.Course.Type,
			CourseSn:    lesson.Course.CourseSn,
			Term:        int32(lesson.Course.Term),
			LessonTotal: lesson.Course.LessonTotal,
			Desc:        lesson.Course.Desc,
			Image:       lesson.Course.Image,
			Teacher:     lesson.Course.Teacher,
			Credit:      lesson.Course.Credit,
			UserId:      lesson.Course.UserId,
		},
		Term:   int32(lesson.Term),
		UserId: lesson.UserId,
	}, nil
}
