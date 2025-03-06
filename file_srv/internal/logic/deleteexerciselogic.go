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

type DeleteExerciseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteExerciseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteExerciseLogic {
	return &DeleteExerciseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteExerciseLogic) DeleteExercise(in *proto.DeleteExerciseRequest) (*proto.Empty, error) {
	result := l.svcCtx.DB.Model(&model.Exercise{}).Delete(&model.Exercise{BaseModel: model.BaseModel{ID: in.Id}})

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的课时ID")
	}

	return &proto.Empty{}, nil
}
