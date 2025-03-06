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

type DeleteCourseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCourseLogic {
	return &DeleteCourseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCourseLogic) DeleteCourse(in *proto.DeleteCourseRequest) (*proto.Empty, error) {
	result := l.svcCtx.DB.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: in.Id}}).Delete(&model.Course{})
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效课程ID")
	}
	return &proto.Empty{}, nil
}
