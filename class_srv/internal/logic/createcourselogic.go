package logic

import (
	"LearningGuide/class_srv/internal/model"
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	proto "LearningGuide/class_srv/.ClassProto"
	"LearningGuide/class_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCourseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCourseLogic {
	return &CreateCourseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCourseLogic) CreateCourse(in *proto.CreateCourseRequest) (*proto.CreateCourseResponse, error) {
	course := model.Course{
		UserId:      in.UserId,
		Name:        in.Name,
		Type:        in.Type,
		Image:       in.Image,
		Teacher:     in.Teacher,
		Credit:      in.Credit,
		LessonTotal: 0,
		Desc:        in.Desc,
		Term:        int(in.Term),
	}

	tx := l.svcCtx.DB.Begin()

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
