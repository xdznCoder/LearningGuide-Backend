package handler

import (
	"LearningGuide/class_srv/global"
	"LearningGuide/class_srv/model"
	proto "LearningGuide/class_srv/proto/.ClassProto"
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"
)

type ClassServer struct {
	proto.UnimplementedClassServer
}

func (c ClassServer) GetCourseList(ctx context.Context, req *proto.CourseFilterRequest) (*proto.CourseListResponse, error) {
	var CourseList []model.Course
	var count int64

	filter := global.DB.Model(&model.Course{})

	if req.UserId != 0 {
		filter = filter.Where(&model.Course{UserId: req.UserId})
	}

	if req.Teacher != "" {
		filter = filter.Where("teacher LIKE ?", "%"+req.Teacher+"%")
	}

	if req.Type != "" {
		filter = filter.Where(&model.Course{Type: req.Type})
	}

	if req.Term != 0 {
		filter = filter.Where(&model.Course{Term: int(req.Term)})
	}

	if req.Name != "" {
		filter = filter.Where("name LIKE ?", "%"+req.Name+"%")
	}

	filter.Count(&count)

	result := filter.Scopes(Paginate(int(req.PageNum), int(req.PageSize))).Model(&model.Course{}).Find(&CourseList)
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

func (c ClassServer) CreateCourse(ctx context.Context, req *proto.CreateCourseRequest) (*proto.CreateCourseResponse, error) {
	course := model.Course{
		UserId:      req.UserId,
		Name:        req.Name,
		Type:        req.Type,
		Image:       req.Image,
		Teacher:     req.Teacher,
		Credit:      req.Credit,
		LessonTotal: 0,
		Desc:        req.Desc,
		Term:        int(req.Term),
	}

	tx := global.DB.Begin()

	result := tx.Model(&model.Course{}).Create(&course)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建课程失败: %v", result.Error)
	}

	result = tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: course.ID}}).Update("course_sn", newCourseSn(course.ID))
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建课程失败: %v", result.Error)
	}

	tx.Commit()

	return &proto.CreateCourseResponse{Id: course.ID}, nil
}

func (c ClassServer) GetCourseDetail(ctx context.Context, req *proto.CourseDetailRequest) (*proto.CourseInfoResponse, error) {
	var course model.Course
	result := global.DB.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: req.Id}}).Find(&course)

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

func (c ClassServer) UpdateCourse(ctx context.Context, req *proto.UpdateCourseRequest) (*proto.Empty, error) {
	course := model.Course{
		Name:    req.Name,
		Image:   req.Image,
		Teacher: req.Teacher,
		Credit:  req.Credit,
		Desc:    req.Desc,
	}

	result := global.DB.Where(&model.Course{BaseModel: model.BaseModel{ID: req.Id}}).Updates(&course)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效课程ID")
	}
	return &proto.Empty{}, nil
}

func (c ClassServer) DeleteCourse(ctx context.Context, req *proto.DeleteCourseRequest) (*proto.Empty, error) {
	result := global.DB.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: req.Id}}).Delete(&model.Course{})
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效课程ID")
	}
	return &proto.Empty{}, nil
}

func Paginate(pageNum, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNum == 0 {
			pageNum = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func newCourseSn(Id int32) string {
	CreateTime := time.Now()
	return fmt.Sprintf("%d%d%d%d%d%d%d%d",
		CreateTime.Year(),
		CreateTime.Month(),
		CreateTime.Day(),
		CreateTime.Hour(),
		CreateTime.Minute(),
		CreateTime.Second(),
		Id,
		random.RandInt(10, 99),
	)
}
