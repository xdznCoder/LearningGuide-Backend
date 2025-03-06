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

type GetCourseDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCourseDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCourseDetailLogic {
	return &GetCourseDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCourseDetailLogic) GetCourseDetail(in *proto.CourseDetailRequest) (*proto.CourseInfoResponse, error) {
	var course model.Course
	result := l.svcCtx.DB.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: in.Id}}).Find(&course)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效课程ID")
	}

	return &proto.CourseInfoResponse{
		Id:          course.ID,
		UserId:      course.UserId,
		Name:        course.Name,
		Type:        course.Type,
		CourseSn:    course.CourseSn,
		Term:        int32(course.Term),
		LessonTotal: course.LessonTotal,
		Desc:        course.Desc,
		Image:       course.Image,
		Teacher:     course.Teacher,
		Credit:      course.Credit,
	}, nil
}
