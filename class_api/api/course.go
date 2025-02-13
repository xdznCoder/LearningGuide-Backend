package api

import (
	"LearningGuide/class_api/forms"
	"LearningGuide/class_api/global"
	"LearningGuide/class_api/middlewares"
	proto "LearningGuide/class_api/proto/.ClassProto"
	"LearningGuide/class_api/utils"
	"LearningGuide/class_api/validator"
	"context"
	lancet "github.com/duke-git/lancet/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func GetCourseList(c *gin.Context) {
	ctx := contextWithSpan(c)

	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.DefaultQuery("name", "")
	courseType := c.DefaultQuery("type", "")
	teacher := c.DefaultQuery("teacher", "")
	term, _ := strconv.Atoi(c.DefaultQuery("term", "0"))
	userId, _ := strconv.Atoi(c.DefaultQuery("userId", "0"))

	if term < 0 || term > 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效学期参数",
		})
	}

	resp, err := global.ClassSrvClient.GetCourseList(ctx, &proto.CourseFilterRequest{
		UserId:   int32(userId),
		Name:     name,
		Type:     courseType,
		Term:     int32(term),
		PageNum:  int32(pageNum),
		PageSize: int32(pageSize),
		Teacher:  teacher,
	})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func NewCourse(c *gin.Context) {
	ctx := contextWithSpan(c)

	var courseForm forms.CreateCourseForm

	err := c.ShouldBind(&courseForm)
	if err != nil {
		validator.HandleValidatorError(err, c)
		return
	}

	if !checkIfAuthorized(courseForm.UserId, c) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无权添加课程",
		})
		return
	}

	resp, err := global.ClassSrvClient.CreateCourse(ctx, &proto.CreateCourseRequest{
		Name:    courseForm.Name,
		Type:    courseForm.Type,
		Term:    int32(courseForm.Term),
		Desc:    courseForm.Desc,
		Image:   courseForm.Image,
		Credit:  courseForm.Credit,
		Teacher: courseForm.Teacher,
		UserId:  int32(courseForm.UserId),
	})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func UpdateCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效路径参数",
		})
		return
	}

	ctx := contextWithSpan(c)

	var courseForm forms.UpdateCourseForm

	err = c.ShouldBindJSON(&courseForm)
	if err != nil {
		validator.HandleValidatorError(err, c)
		return
	}

	resp, err := global.ClassSrvClient.GetCourseDetail(ctx, &proto.CourseDetailRequest{Id: int32(id)})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	if !checkIfAuthorized(int(resp.UserId), c) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无权更改",
		})
		return
	}

	if courseForm.Image != "" && !lancet.IsUrl(courseForm.Image) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "image必须为有效的url",
		})
		return
	}

	_, err = global.ClassSrvClient.UpdateCourse(ctx, &proto.UpdateCourseRequest{
		Id:      int32(id),
		Name:    courseForm.Name,
		Desc:    courseForm.Desc,
		Image:   courseForm.Image,
		Teacher: courseForm.Teacher,
		Credit:  courseForm.Credit,
	})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

func DeleteCourse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效路径参数",
		})
		return
	}

	ctx := contextWithSpan(c)

	resp, err := global.ClassSrvClient.GetCourseDetail(ctx, &proto.CourseDetailRequest{Id: int32(id)})

	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	if checkIfAuthorized(int(resp.UserId), c) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "无权更改",
		})
		return
	}

	_, err = global.ClassSrvClient.DeleteCourse(ctx, &proto.DeleteCourseRequest{Id: int32(id)})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func GetCourseDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效路径参数",
		})
		return
	}

	ctx := contextWithSpan(c)

	resp, err := global.ClassSrvClient.GetCourseDetail(ctx, &proto.CourseDetailRequest{Id: int32(id)})
	if err != nil {
		utils.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func checkIfAuthorized(id int, c *gin.Context) bool {
	user, _ := c.Get("claims")
	userId := user.(*middlewares.CustomClaims).ID
	role := user.(*middlewares.CustomClaims).AuthorityId
	if uint(id) == userId {
		return true
	}
	if role != 1 {
		return true
	}

	return false
}

func contextWithSpan(c *gin.Context) context.Context {
	ctx := context.Background()
	span, ok := c.Get("span")

	if !ok {
		zap.S().Info(c.Request.URL.Path + "no tracer injected")
		return ctx
	}

	ctx = opentracing.ContextWithSpan(ctx, span.(opentracing.Span))
	return ctx
}
