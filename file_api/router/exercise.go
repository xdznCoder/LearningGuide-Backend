package router

import (
	"LearningGuide/file_api/api"
	"github.com/OuterCyrex/Gorra/GorraAPI/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func ExerciseRouter(R *gin.RouterGroup) {
	Exercise := R.Group("/exer").Use(middlewares.TracerMiddleware(opentracing.GlobalTracer()))
	{
		Exercise.POST("", api.NewExercise)
		Exercise.GET("/list", api.ExerciseList)
		Exercise.GET("/:id", api.GetExercise)
		Exercise.PUT("/:id", api.UpdateRight)
		Exercise.DELETE("/:id", api.DeleteExercise)
	}
}
