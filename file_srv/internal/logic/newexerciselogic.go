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

type NewExerciseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNewExerciseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewExerciseLogic {
	return &NewExerciseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NewExerciseLogic) NewExercise(in *proto.NewExerciseRequest) (*proto.NewExerciseResponse, error) {
	exercise := model.Exercise{
		CourseId: in.CourseId,
		Question: in.Question,
		SectionA: in.SectionA,
		SectionB: in.SectionB,
		SectionC: in.SectionC,
		SectionD: in.SectionD,
		Answer:   in.Answer,
		Reason:   in.Reason,
	}

	result := l.svcCtx.DB.Model(&model.Exercise{}).Create(&exercise)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &proto.NewExerciseResponse{Id: exercise.ID}, nil
}
