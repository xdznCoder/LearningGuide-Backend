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

type DeleteLessonLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLessonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLessonLogic {
	return &DeleteLessonLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLessonLogic) DeleteLesson(in *proto.DeleteLessonRequest) (*proto.Empty, error) {
	var lesson model.Lesson

	tx := l.svcCtx.DB.Begin()

	result := tx.Model(&model.Lesson{}).Where(&model.Lesson{BaseModel: model.BaseModel{ID: in.Id}}).Find(&lesson)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "无效课时ID")
	}

	var course model.Course
	result = tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: lesson.CourseId}}).Find(&course)

	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	result = tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: lesson.CourseId}}).Update("lesson_total", course.LessonTotal-1)

	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	result = tx.Model(&model.Lesson{}).Delete(&model.Lesson{BaseModel: model.BaseModel{ID: in.Id}})
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "无效课时ID")
	}

	tx.Commit()

	return &proto.Empty{}, nil
}
