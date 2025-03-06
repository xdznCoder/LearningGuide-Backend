package logic

import (
	"LearningGuide/file_srv/internal/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "LearningGuide/file_srv/.FileProto"
	"LearningGuide/file_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateExerciseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateExerciseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateExerciseLogic {
	return &UpdateExerciseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateExerciseLogic) UpdateExercise(in *proto.UpdateExerciseRequest) (*proto.Empty, error) {
	switch in.IsRight {
	case "true":
	case "false":
	default:
		return nil, status.Errorf(codes.InvalidArgument, "无效的is_right参数")
	}

	result := l.svcCtx.DB.Model(&model.Exercise{}).
		Where(&model.Exercise{BaseModel: model.BaseModel{ID: in.Id}}).
		Update("is_right", in.IsRight)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的习题ID")
	}

	return &proto.Empty{}, nil
}
