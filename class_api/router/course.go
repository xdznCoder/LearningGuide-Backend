package router

import (
	"LearningGuide/class_api/api"
	"LearningGuide/class_api/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func InitCourseRouter(Router *gin.RouterGroup) {
	zap.S().Infof("Initialize LearningGuide User Router...")
	UserRouter := Router.Group("/course").Use(middlewares.JWTAuthMiddleware()).Use(middlewares.TracerMiddleware(opentracing.GlobalTracer()))
	{
		UserRouter.GET("list", api.GetCourseList)
		UserRouter.POST("", api.NewCourse)
		UserRouter.PUT("/:id", api.UpdateCourse)
		UserRouter.DELETE("/:id", api.DeleteCourse)
		UserRouter.GET("/:id", api.GetCourseDetail)
	}
}
