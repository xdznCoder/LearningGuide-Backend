package api

import (
	"LearningGuide/file_api/forms"
	"LearningGuide/file_api/global"
	ChatProto "LearningGuide/file_api/proto/.ChatProto"
	FileProto "LearningGuide/file_api/proto/.FileProto"
	"errors"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	handleGrpc "github.com/OuterCyrex/Gorra/GorraAPI/resp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func ExerciseList(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	pageSize, err1 := strconv.Atoi(c.DefaultQuery("pageSize", "0"))
	pageNum, err2 := strconv.Atoi(c.DefaultQuery("pageNum", "0"))
	courseId, err3 := strconv.Atoi(c.DefaultQuery("course_id", "0"))
	if errors.Join(err1, err2, err3) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的查询参数",
		})
		return
	}
	isRight := c.DefaultQuery("is_right", "")
	question := c.DefaultQuery("question", "")

	switch isRight {
	case "":
	case "true":
	case "false":
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的is_right参数",
		})
		return
	}

	resp, err := global.FileSrvClient.ExerciseList(ctx, &FileProto.ExerciseListRequest{
		CourseId: int32(courseId),
		Question: question,
		IsRight:  isRight,
		PageNum:  int32(pageNum),
		PageSize: int32(pageSize),
	})
	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func NewExercise(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	var newExerForms forms.NewExerciseForms

	err := c.ShouldBindJSON(&newExerForms)

	if err != nil {
		handleGrpc.HandleValidatorError(err, c)
		return
	}

	stream, err := global.ChatSrvClient.SendStreamMessage(ctx, &ChatProto.UserMessage{
		CourseID:     newExerForms.CourseId,
		TemplateType: int32(TemplateTypeExerciseGenerate),
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	result, err := ToString(stream)

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	exer, err := transResultToExercise(result)

	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请重试",
		})
		zap.S().Errorf("transmit glm result to exercise failed: %v", err)
		return
	}

	resp, err := global.FileSrvClient.NewExercise(ctx, &FileProto.NewExerciseRequest{
		CourseId: newExerForms.CourseId,
		Question: exer.Question,
		SectionA: exer.Sections.A,
		SectionB: exer.Sections.B,
		SectionC: exer.Sections.C,
		SectionD: exer.Sections.D,
		Answer:   exer.Answer,
		Reason:   exer.Reason,
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func UpdateRight(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的路径参数",
		})
		return
	}

	updateForm := forms.UpdateExerciseForms{}

	err = c.ShouldBindJSON(&updateForm)

	if err != nil {
		handleGrpc.HandleValidatorError(err, c)
		return
	}

	_, err = global.FileSrvClient.UpdateExercise(ctx, &FileProto.UpdateExerciseRequest{
		IsRight: updateForm.IsRight,
		Id:      int32(id),
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

func GetExercise(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的路径参数",
		})
		return
	}

	resp, err := global.FileSrvClient.GetExerciseDetail(ctx, &FileProto.ExerciseDetailRequest{Id: int32(id)})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func DeleteExercise(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的路径参数",
		})
		return
	}

	_, err = global.FileSrvClient.DeleteExercise(ctx, &FileProto.DeleteExerciseRequest{Id: int32(id)})
	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}
