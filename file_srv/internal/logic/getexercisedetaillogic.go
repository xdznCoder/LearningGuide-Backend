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

type GetExerciseDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetExerciseDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetExerciseDetailLogic {
	return &GetExerciseDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetExerciseDetailLogic) GetExerciseDetail(in *proto.ExerciseDetailRequest) (*proto.ExerciseInfoResponse, error) {
	var exer model.Exercise

	result := l.svcCtx.DB.Model(&model.Exercise{}).Where(&model.Exercise{BaseModel: model.BaseModel{ID: in.Id}}).Find(&exer)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的习题ID")
	}

	return &proto.ExerciseInfoResponse{
		CourseId: exer.CourseId,
		Question: exer.Question,
		SectionA: exer.SectionA,
		SectionB: exer.SectionB,
		SectionC: exer.SectionC,
		SectionD: exer.SectionD,
		Answer:   exer.Answer,
		Reason:   exer.Reason,
		IsRight:  exer.IsRight,
		Id:       exer.ID,
	}, nil
}
