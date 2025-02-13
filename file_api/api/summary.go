package api

import (
	"LearningGuide/file_api/forms"
	"LearningGuide/file_api/global"
	FileProto "LearningGuide/file_api/proto/.FileProto"
	"fmt"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	handleGrpc "github.com/OuterCyrex/Gorra/GorraAPI/resp"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func NewSummary(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	var newForm forms.NewSummaryForms
	err := c.ShouldBindJSON(&newForm)

	if err != nil {
		handleGrpc.HandleValidatorError(err, c)
		return
	}

	year, week := time.Now().ISOWeek()

	var isoWeek string

	if week > 1 {
		isoWeek = fmt.Sprintf("%04d%02d", year, week-1)
	} else {
		isoWeek = fmt.Sprintf("%04d%02d", year, week)
	}

	if newForm.ISOWeek != "" {
		if !isISOWeekFormat(newForm.ISOWeek) {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "无效的iso_week参数",
			})
			return
		}

		isoWeek = newForm.ISOWeek
	}

	resp, err := global.FileSrvClient.NewSummary(ctx, &FileProto.NewSummaryRequest{
		CourseID: newForm.CourseId,
		ISOWeek:  isoWeek,
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetSummary(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的路径参数",
		})
		return
	}

	resp, err := global.FileSrvClient.GetSummary(ctx, &FileProto.GetSummaryRequest{Id: int32(id)})
	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, grpcToHttpResponse(resp))
}

func SummaryList(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	courseId, err := strconv.Atoi(c.DefaultQuery("course_id", "0"))
	pageNum, err := strconv.Atoi(c.DefaultQuery("pageNum", "0"))
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	year := c.DefaultQuery("year", "")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的查询参数",
		})
		return
	}

	resp, err := global.FileSrvClient.SummaryList(ctx, &FileProto.SummaryListRequest{
		Year:     year,
		CourseId: int32(courseId),
		PageNum:  int32(pageNum),
		PageSize: int32(pageSize),
	})
	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	var data []gin.H

	for _, v := range resp.Data {
		data = append(data, grpcToHttpResponse(v))
	}

	c.JSON(http.StatusOK, gin.H{
		"total": resp.Total,
		"data":  data,
	})
}

func grpcToHttpResponse(resp *FileProto.SummaryInfoResponse) gin.H {
	var (
		ExerciseTotal int32
		TrueTotal     float32
	)
	for i := 0; i < 7; i++ {
		ExerciseTotal += resp.ExerciseDone[i]
		TrueTotal += resp.AccuracyRate[i] * float32(resp.ExerciseDone[i])
	}

	var AccuracyTotal float32

	if ExerciseTotal != 0 {
		AccuracyTotal = TrueTotal / float32(ExerciseTotal)
	} else {
		AccuracyTotal = 0
	}

	return gin.H{
		"id":                  resp.Id,
		"week_id":             resp.WeekID,
		"course_id":           resp.CourseID,
		"accuracy_rate_array": resp.AccuracyRate,
		"exercise_done_array": resp.ExerciseDone,
		"exercise_done_total": ExerciseTotal,
		"accuracy_rate_total": AccuracyTotal,
		"noun_total":          resp.NounNum,
		"session_total":       resp.SessionNum,
		"message_total":       resp.MessageNum,
	}
}

func isISOWeekFormat(s string) bool {
	regex := `^\d{4}\d{2}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(s)
}
