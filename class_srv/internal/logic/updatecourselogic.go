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

type UpdateCourseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCourseLogic {
	return &UpdateCourseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCourseLogic) UpdateCourse(in *proto.UpdateCourseRequest) (*proto.Empty, error) {
	course := model.Course{
		Name:    in.Name,
		Image:   in.Image,
		Teacher: in.Teacher,
		Credit:  in.Credit,
		Desc:    in.Desc,
	}

	result := l.svcCtx.DB.Where(&model.Course{BaseModel: model.BaseModel{ID: in.Id}}).Updates(&course)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效课程ID")
	}
	return &proto.Empty{}, nil
}
