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

type SummaryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSummaryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SummaryListLogic {
	return &SummaryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SummaryListLogic) SummaryList(req *proto.SummaryListRequest) (*proto.SummaryListResponse, error) {
	var summaries []model.Summary
	var count int64

	l.svcCtx.DB.Model(&model.Summary{}).Where(&model.Summary{CourseID: req.CourseId}).
		Where("week_id LIKE ?", req.Year+"%").
		Count(&count)

	result := l.svcCtx.DB.Model(&model.Summary{}).Scopes(model.Paginate(int(req.PageNum), int(req.PageSize))).Where(&model.Summary{
		CourseID: req.CourseId,
	}).Where("week_id LIKE ?", req.Year+"%").Find(&summaries)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	var respList []*proto.SummaryInfoResponse

	for _, v := range summaries {
		exerciseArray := make([]int32, 7)
		accuracyArray := make([]float32, 7)

		exerTmp := strings.Split(v.ExerciseDone, ",")
		accuTmp := strings.Split(v.AccuracyRate, ",")

		for i := 0; i < 7; i++ {
			tmpInt, _ := strconv.Atoi(exerTmp[i])
			exerciseArray[i] = int32(tmpInt)

			tmpFloat, _ := strconv.ParseFloat(accuTmp[i], 32)
			accuracyArray[i] = float32(tmpFloat)
		}

		respList = append(respList, &proto.SummaryInfoResponse{
			Id:           v.ID,
			WeekID:       v.WeekID,
			CourseID:     v.CourseID,
			ExerciseDone: exerciseArray,
			AccuracyRate: accuracyArray,
			SessionNum:   v.SessionNum,
			MessageNum:   v.MessageNum,
			NounNum:      v.NounNum,
		})
	}

	return &proto.SummaryListResponse{
		Total: int32(count),
		Data:  respList,
	}, nil
}
