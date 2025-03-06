package logic

import (
	"LearningGuide/class_srv/internal/model"
	"context"
	"github.com/duke-git/lancet/v2/slice"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "LearningGuide/class_srv/.ClassProto"
	"LearningGuide/class_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLessonInBatchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLessonInBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLessonInBatchLogic {
	return &DeleteLessonInBatchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLessonInBatchLogic) DeleteLessonInBatch(in *proto.DeleteLessonInBatchRequest) (*proto.Empty, error) {
	in.Ids = slice.Unique(in.Ids)

	tx := l.svcCtx.DB.Begin()

	var course model.Course
	result := tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: in.CourseId}}).Find(&course)

	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "无效课程ID")
	}

	result = tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: in.CourseId}}).Update("lesson_total", int(course.LessonTotal)-len(in.Ids))

	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	var conditions [][]interface{}
	for _, id := range in.Ids {
		conditions = append(conditions, []interface{}{id, in.UserId, in.CourseId})
	}

	result = tx.Model(&model.Lesson{}).Where("(id, user_id, course_id) IN (?)", conditions).Delete(&model.Lesson{})
	if result.RowsAffected < int64(len(in.Ids)) {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "存在无效课时ID")
	}

	tx.Commit()
	return &proto.Empty{}, nil
}
