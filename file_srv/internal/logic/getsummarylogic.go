package logic

import (
	"LearningGuide/file_srv/internal/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"

	proto "LearningGuide/file_srv/.FileProto"
	"LearningGuide/file_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSummaryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSummaryLogic {
	return &GetSummaryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSummaryLogic) GetSummary(req *proto.GetSummaryRequest) (*proto.SummaryInfoResponse, error) {
	var summary model.Summary

	result := l.svcCtx.DB.Model(&model.Summary{}).Where(&model.Summary{BaseModel: model.BaseModel{ID: req.Id}}).Find(&summary)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效总结ID")
	}

	exerciseArray := make([]int32, 7)
	accuracyArray := make([]float32, 7)

	exerTmp := strings.Split(summary.ExerciseDone, ",")
	accuTmp := strings.Split(summary.AccuracyRate, ",")

	for i := 0; i < 7; i++ {
		tmpInt, _ := strconv.Atoi(exerTmp[i])
		exerciseArray[i] = int32(tmpInt)

		tmpFloat, _ := strconv.ParseFloat(accuTmp[i], 32)
		accuracyArray[i] = float32(tmpFloat)
	}

	return &proto.SummaryInfoResponse{
		Id:           summary.ID,
		WeekID:       summary.WeekID,
		CourseID:     summary.CourseID,
		ExerciseDone: exerciseArray,
		AccuracyRate: accuracyArray,
		SessionNum:   summary.SessionNum,
		MessageNum:   summary.MessageNum,
		NounNum:      summary.NounNum,
	}, nil
}
