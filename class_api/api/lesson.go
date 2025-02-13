package api

import (
	"LearningGuide/class_api/forms"
	"LearningGuide/class_api/global"
	proto "LearningGuide/class_api/proto/.ClassProto"
	"LearningGuide/class_api/utils"
	"LearningGuide/class_api/validator"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetLessonList(c *gin.Context) {
	ctx := contextWithSpan(c)

	cid, err := strconv.Atoi(c.DefaultQuery("course_id", "0"))
	weekNum, err := strconv.Atoi(c.DefaultQuery("week_num", "0"))
	dayOfWeek, err := strconv.Atoi(c.DefaultQuery("day_of_week", "0"))
	term, err := strconv.Atoi(c.DefaultQuery("term", "0"))
	userId, err := strconv.Atoi(c.DefaultQuery("userId", "0"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效参数",
		})
		return
	}

	if term < 0 || term > 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效学期参数，应处于[1,8]区间内",
		})
		return
	}

	resp, err := global.ClassSrvClient.GetLessonList(ctx, &proto.LessonFilterRequest{
		CourseId:  int32(cid),
		WeekNum:   int32(weekNum),
		DayOfWeek: int32(dayOfWeek),
		Term:      int32(term),
		UserId:    int32(userId),
	})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func NewLesson(c *gin.Context) {
	ctx := contextWithSpan(c)

	lessonForm := forms.CreateLessonForm{}

	err := c.ShouldBindJSON(&lessonForm)
	if err != nil {
		validator.HandleValidatorError(err, c)
		return
	}

	if !checkIfAuthorized(lessonForm.UserId, c) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无权添加",
		})
		return
	}

	resp, err := global.ClassSrvClient.CreateLesson(ctx, &proto.CreateLessonRequest{
		CourseId:  int32(lessonForm.CourseId),
		WeekNum:   int32(lessonForm.WeekNum),
		DayOfWeek: int32(lessonForm.DayOfWeek),
		LessonNum: int32(lessonForm.LessonNum),
		Begin:     lessonForm.Begin,
		End:       lessonForm.End,
		UserId:    int32(lessonForm.UserId),
	})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func UpdateLesson(c *gin.Context) {
	ctx := contextWithSpan(c)

	lessonForm := forms.UpdateLessonForm{}

	err := c.ShouldBindJSON(&lessonForm)
	if err != nil {
		validator.HandleValidatorError(err, c)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效路径参数",
		})
		return
	}

	resp, err := global.ClassSrvClient.GetLessonDetail(ctx, &proto.LessonDetailRequest{Id: int32(id)})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	if !checkIfAuthorized(int(resp.UserId), c) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无权删除",
		})
		return
	}

	_, err = global.ClassSrvClient.UpdateLesson(ctx, &proto.UpdateLessonRequest{
		Id:    int32(id),
		Begin: lessonForm.Begin,
		End:   lessonForm.End,
	})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更改成功",
	})
}

func DeleteLesson(c *gin.Context) {
	ctx := contextWithSpan(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效路径参数",
		})
		return
	}

	resp, err := global.ClassSrvClient.GetLessonDetail(ctx, &proto.LessonDetailRequest{Id: int32(id)})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	if !checkIfAuthorized(int(resp.UserId), c) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无权删除",
		})
		return
	}

	_, err = global.ClassSrvClient.DeleteLesson(ctx, &proto.DeleteLessonRequest{
		Id: int32(id),
	})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func NewLessonInBatch(c *gin.Context) {
	ctx := contextWithSpan(c)

	form := forms.CreateLessonFormInBatch{}

	err := c.ShouldBindJSON(&form)
	if err != nil {
		validator.HandleValidatorError(err, c)
		return
	}

	if !checkIfAuthorized(form.UserId, c) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无权添加",
		})
		return
	}

	if form.EndWeek < form.BeginWeek {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效始末周数",
		})
		return
	}

	resp, err := global.ClassSrvClient.CreateLessonInBatch(ctx, &proto.CreateLessonBatchRequest{
		CourseId:  int32(form.CourseId),
		UserId:    int32(form.UserId),
		BeginWeek: int32(form.BeginWeek),
		EndWeek:   int32(form.EndWeek),
		DayOfWeek: int32(form.DayOfWeek),
		LessonNum: int32(form.LessonNum),
		Begin:     form.Begin,
		End:       form.End,
	})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetLessonDetail(c *gin.Context) {
	ctx := contextWithSpan(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效路径参数",
		})
		return
	}

	resp, err := global.ClassSrvClient.GetLessonDetail(ctx, &proto.LessonDetailRequest{Id: int32(id)})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	if !checkIfAuthorized(int(resp.UserId), c) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无权访问",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func DeleteLessonInBatch(c *gin.Context) {
	ctx := contextWithSpan(c)

	deleteForm := forms.DeleteLessonFormInBatch{}

	err := c.ShouldBindJSON(&deleteForm)
	if err != nil {
		validator.HandleValidatorError(err, c)
		return
	}

	if !checkIfAuthorized(deleteForm.UserId, c) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无权添加",
		})
		return
	}

	_, err = global.ClassSrvClient.DeleteLessonInBatch(ctx, &proto.DeleteLessonInBatchRequest{
		UserId:   int32(deleteForm.UserId),
		CourseId: int32(deleteForm.CourseId),
		Ids:      deleteForm.Ids,
	})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "全部删除成功",
	})
}
