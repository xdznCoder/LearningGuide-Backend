package router

import (
	"LearningGuide/file_api/api"
	"github.com/OuterCyrex/Gorra/GorraAPI/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func SummaryRouter(R *gin.RouterGroup) {
	Summary := R.Group("/sum").Use(middlewares.TracerMiddleware(opentracing.GlobalTracer()))
	{
		Summary.POST("", api.NewSummary)
		Summary.GET("/list", api.SummaryList)
		Summary.GET("/:id", api.GetSummary)
	}
}
