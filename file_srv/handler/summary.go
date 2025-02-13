package handler

import (
	"LearningGuide/file_srv/global"
	"LearningGuide/file_srv/model"
	FileProto "LearningGuide/file_srv/proto/.FileProto"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (f FileServer) NewSummary(ctx context.Context, req *FileProto.NewSummaryRequest) (*FileProto.NewSummaryResponse, error) {
	if !isISOWeekFormat(req.ISOWeek) {
		return nil, status.Errorf(codes.InvalidArgument, "无效ISOWeek格式")
	}

	year, err := strconv.Atoi(req.ISOWeek[:4])
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

	global.DB.Model(&model.Session{}).Where("course_id = ?", req.CourseID).
		Where("update_time >= ? AND update_time < ?", firstDayOfWeek, lastDayOfWeek).Find(&sessions)

	if len(sessions) > 0 {
		for _, v := range sessions {
			sessionId = append(sessionId, v.ID)
		}

		global.DB.Model(&model.Message{}).Where("session_id IN (?)", sessionId).
			Where("update_time >= ? AND update_time < ?", firstDayOfWeek, lastDayOfWeek).Count(&messageNum)
	}

	global.DB.Model(&model.Noun{}).Where("course_id = ?", req.CourseID).
		Where("update_time >= ? AND update_time < ?", firstDayOfWeek, lastDayOfWeek).Count(&nounNum)

	for i := 0; i < 7; i++ {
		day := firstDayOfWeek.AddDate(0, 0, i)
		nextDay := day.AddDate(0, 0, 1)

		var (
			trueNum  int64
			falseNum int64
		)

		global.DB.Model(&model.Exercise{}).
			Where("course_id = ?", req.CourseID).
			Where("update_time >= ? AND update_time < ?", day, nextDay).
			Where(&model.Exercise{IsRight: "true"}).
			Count(&trueNum)

		global.DB.Model(&model.Exercise{}).
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

	result := global.DB.Model(&model.Summary{}).
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

	result = global.DB.Save(&summary)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &FileProto.NewSummaryResponse{Id: summary.ID}, nil
}

func (f FileServer) SummaryList(ctx context.Context, req *FileProto.SummaryListRequest) (*FileProto.SummaryListResponse, error) {
	var summaries []model.Summary
	var count int64

	global.DB.Model(&model.Summary{}).Where(&model.Summary{CourseID: req.CourseId}).
		Where("week_id LIKE ?", req.Year+"%").
		Count(&count)

	result := global.DB.Model(&model.Summary{}).Scopes(Paginate(int(req.PageNum), int(req.PageSize))).Where(&model.Summary{
		CourseID: req.CourseId,
	}).Where("week_id LIKE ?", req.Year+"%").Find(&summaries)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	var respList []*FileProto.SummaryInfoResponse

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

		respList = append(respList, &FileProto.SummaryInfoResponse{
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

	return &FileProto.SummaryListResponse{
		Total: int32(count),
		Data:  respList,
	}, nil
}

func (f FileServer) GetSummary(ctx context.Context, req *FileProto.GetSummaryRequest) (*FileProto.SummaryInfoResponse, error) {
	var summary model.Summary

	result := global.DB.Model(&model.Summary{}).Where(&model.Summary{BaseModel: model.BaseModel{ID: req.Id}}).Find(&summary)

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

	return &FileProto.SummaryInfoResponse{
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

func isISOWeekFormat(s string) bool {
	regex := `^\d{4}\d{2}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(s)
}
