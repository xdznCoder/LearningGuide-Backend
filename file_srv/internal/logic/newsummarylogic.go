package logic

import (
	"LearningGuide/file_srv/internal/model"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"strconv"
	"strings"
	"time"

	proto "LearningGuide/file_srv/.FileProto"
	"LearningGuide/file_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type NewSummaryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNewSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewSummaryLogic {
	return &NewSummaryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NewSummaryLogic) NewSummary(req *proto.NewSummaryRequest) (*proto.NewSummaryResponse, error) {
	if !isISOWeekFormat(req.ISOWeek) {
		return nil, status.Errorf(codes.InvalidArgument, "无效ISOWeek格式")
	}

	year, err := strconv.Atoi(req.ISOWeek[:4])
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "无效ISOWeek格式")
	}

	week, err := strconv.Atoi(req.ISOWeek[4:])

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "无效ISOWeek格式")
	}

	startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	firstMonday := startOfYear.AddDate(0, 0, -int(startOfYear.Weekday())+1)

	firstDayOfWeek := firstMonday.AddDate(0, 0, (week-1)*7)
	lastDayOfWeek := firstDayOfWeek.AddDate(0, 0, 6)

	exerciseArray := make([]string, 7)
	accuracyArray := make([]string, 7)

	var (
		sessionNum int64
		messageNum int64
		nounNum    int64
	)

	var sessions []model.Session
	var sessionId []int32

	l.svcCtx.DB.Model(&model.Session{}).Where("course_id = ?", req.CourseID).
		Where("update_time >= ? AND update_time < ?", firstDayOfWeek, lastDayOfWeek).Find(&sessions)

	if len(sessions) > 0 {
		for _, v := range sessions {
			sessionId = append(sessionId, v.ID)
		}

		l.svcCtx.DB.Model(&model.Message{}).Where("session_id IN (?)", sessionId).
			Where("update_time >= ? AND update_time < ?", firstDayOfWeek, lastDayOfWeek).Count(&messageNum)
	}

	l.svcCtx.DB.Model(&model.Noun{}).Where("course_id = ?", req.CourseID).
		Where("update_time >= ? AND update_time < ?", firstDayOfWeek, lastDayOfWeek).Count(&nounNum)

	for i := 0; i < 7; i++ {
		day := firstDayOfWeek.AddDate(0, 0, i)
		nextDay := day.AddDate(0, 0, 1)

		var (
			trueNum  int64
			falseNum int64
		)

		l.svcCtx.DB.Model(&model.Exercise{}).
			Where("course_id = ?", req.CourseID).
			Where("update_time >= ? AND update_time < ?", day, nextDay).
			Where(&model.Exercise{IsRight: "true"}).
			Count(&trueNum)

		l.svcCtx.DB.Model(&model.Exercise{}).
			Where("course_id = ?", req.CourseID).
			Where("update_time >= ? AND update_time < ?", day, nextDay).
			Where(&model.Exercise{IsRight: "false"}).
			Count(&falseNum)

		exerciseArray[i] = strconv.Itoa(int(trueNum + falseNum))
		if trueNum+falseNum != 0 {
			accuracyArray[i] = fmt.Sprintf("%2f", float32(trueNum)/float32(trueNum+falseNum))
		} else {
			accuracyArray[i] = "0"
		}
	}

	var former model.Summary

	result := l.svcCtx.DB.Model(&model.Summary{}).
		Where(&model.Summary{WeekID: fmt.Sprintf("%04d%02d", year, week), CourseID: req.CourseID}).
		Find(&former)

	summary := model.Summary{
		WeekID:       req.ISOWeek,
		CourseID:     req.CourseID,
		ExerciseDone: strings.Join(exerciseArray, ","),
		AccuracyRate: strings.Join(accuracyArray, ","),
		SessionNum:   int32(sessionNum),
		MessageNum:   int32(messageNum),
		NounNum:      int32(nounNum),
	}

	if result.RowsAffected != 0 {
		summary.ID = former.ID
		summary.CreatedAt = former.CreatedAt
	}

	result = l.svcCtx.DB.Save(&summary)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &proto.NewSummaryResponse{Id: summary.ID}, nil
}

func isISOWeekFormat(s string) bool {
	regex := `^\d{4}\d{2}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(s)
}
