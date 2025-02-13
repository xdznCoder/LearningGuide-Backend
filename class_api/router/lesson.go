package router

import (
	"LearningGuide/class_api/api"
	"LearningGuide/class_api/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func InitLessonRouter(Router *gin.RouterGroup) {
	zap.S().Infof("Initialize LearningGuide Lesson Router...")
	LessonRouter := Router.Group("/lesson").Use(middlewares.JWTAuthMiddleware()).Use(middlewares.TracerMiddleware(opentracing.GlobalTracer()))
	{
		LessonRouter.GET("list", api.GetLessonList)
		LessonRouter.POST("", api.NewLesson)
		LessonRouter.PUT("/:id", api.UpdateLesson)
		LessonRouter.DELETE("/:id", api.DeleteLesson)
		LessonRouter.POST("/batch", api.NewLessonInBatch)
		LessonRouter.GET("/:id", api.GetLessonDetail)
		LessonRouter.DELETE("/batch", api.DeleteLessonInBatch)
	}
}
