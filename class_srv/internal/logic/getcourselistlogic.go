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

type GetCourseListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCourseListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCourseListLogic {
	return &GetCourseListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCourseListLogic) GetCourseList(in *proto.CourseFilterRequest) (*proto.CourseListResponse, error) {
	var CourseList []model.Course
	var count int64

	filter := l.svcCtx.DB.Model(&model.Course{})

	if in.UserId != 0 {
		filter = filter.Where(&model.Course{UserId: in.UserId})
	}

	if in.Teacher != "" {
		filter = filter.Where("teacher LIKE ?", "%"+in.Teacher+"%")
	}

	if in.Type != "" {
		filter = filter.Where(&model.Course{Type: in.Type})
	}

	if in.Term != 0 {
		filter = filter.Where(&model.Course{Term: int(in.Term)})
	}

	if in.Name != "" {
		filter = filter.Where("name LIKE ?", "%"+in.Name+"%")
	}

	filter.Count(&count)

	result := filter.Scopes(model.Paginate(int(in.PageNum), int(in.PageSize))).Model(&model.Course{}).Find(&CourseList)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "查询课程列表出错：%v", result.Error)
	}

	RespList := make([]*proto.CourseInfoResponse, 0)

	for _, course := range CourseList {
		RespList = append(RespList, &proto.CourseInfoResponse{
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
		})
	}

	return &proto.CourseListResponse{
		Total: int32(count),
		Data:  RespList,
	}, nil
}
