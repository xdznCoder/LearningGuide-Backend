package api

import (
	"LearningGuide/file_api/forms"
	"LearningGuide/file_api/global"
	FileProto "LearningGuide/file_api/proto/.FileProto"
	"encoding/json"
	"errors"
	"github.com/OuterCyrex/ChatGLM_sdk"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	handleGrpc "github.com/OuterCyrex/Gorra/GorraAPI/resp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"strconv"
)

func ExerciseList(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "0"))
	pageNum, err := strconv.Atoi(c.DefaultQuery("pageNum", "0"))
	courseId, err := strconv.Atoi(c.DefaultQuery("course_id", "0"))
	isRight := c.DefaultQuery("is_right", "")
	question := c.DefaultQuery("question", "")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的查询参数",
		})
		return
	}

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

func GenerateExercise(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	exer := forms.GenerateExerciseForms{}

	err := c.ShouldBindJSON(&exer)
	if err != nil {
		handleGrpc.HandleValidatorError(err, c)
		return
	}

	if len(exer.FileIds) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "文件最多可选10个",
		})
		return
	}

	var messageCtx *ChatGLM_sdk.MessageContext

	if len(exer.FileIds) == 0 {
		var ids []int32

		resp, err := global.FileSrvClient.FileList(ctx, &FileProto.FileFilterRequest{
			PageNum:  0,
			PageSize: 10,
			CourseId: exer.CourseId,
		})
		if err != nil {
			handleGrpc.HandleGrpcErrorToHttp(err, c)
			return
		}

		for _, v := range resp.Data {
			ids = append(ids, v.Id)
		}

		messageCtx, err = getFileContext(ctx, ids)
		if err != nil {
			handleGrpc.HandleGrpcErrorToHttp(err, c)
			return
		}
	} else {
		messageCtx, err = getFileContext(ctx, exer.FileIds)
		if err != nil {
			handleGrpc.HandleGrpcErrorToHttp(err, c)
			return
		}
	}

	client := ChatGLM_sdk.NewClient(global.ServerConfig.ChatGLM.AccessKey)

	id, err := client.SendAsync(messageCtx, getExercisePrompt())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务器内部错误",
		})
		zap.S().Errorf("get async id from glm failed: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"glm_result_id": id,
	})
}

func NewExercise(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	var newExerForms forms.NewExerciseForms

	err := c.ShouldBindJSON(&newExerForms)

	if err != nil {
		handleGrpc.HandleValidatorError(err, c)
		return
	}

	client := ChatGLM_sdk.NewClient(global.ServerConfig.ChatGLM.AccessKey)

	result := client.GetAsyncMessage(ChatGLM_sdk.NewContext(), newExerForms.ResultId)

	if errors.Is(result.Error, ChatGLM_sdk.ErrResultProcessing) {
		c.JSON(http.StatusAccepted, gin.H{
			"msg": "GLM正在生成中",
		})
		return
	} else if result.Error != nil || len(result.Message) <= 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务器内部错误",
		})
		zap.S().Errorf("get result from glm failed: %v", result.Error)
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

func transResultToExercise(result ChatGLM_sdk.Result) (exercise, error) {
	pattern := "```json\\s*({[\\s\\S]*?})\\s*```"

	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(result.Message[0].Content)

	var output string

	if matches == nil {
		output = result.Message[0].Content
	} else {
		output = matches[1]
	}

	var question exercise

	err := json.Unmarshal([]byte(output), &question)
	if err != nil {
		return exercise{}, err
	}

	return question, nil
}

type exercise struct {
	Question string      `json:"question"`
	Sections SectionsSet `json:"sections"`
	Answer   string      `json:"answer"`
	Reason   string      `json:"reason"`
}

type SectionsSet struct {
	A string `json:"A"`
	B string `json:"B"`
	C string `json:"C"`
	D string `json:"D"`
}

func getExercisePrompt() string {
	var template = `{
  "question": "问题的内容",
  "sections": {
    "A": "a选项",
    "B": "b选项",
    "C": "c选项",
    "D": "d选项"
  },
  "answer": "本题目的答案",
  "reason": "选择该答案的原因"
}`

	return "请根据所给的文件内容出一个相关练习题，要求题目以该JSON格式返回，JSON格式为：" + template + "不要回复任何多余的话，只需要JSON内容，不要出代码题，答案一定要正确"
}
