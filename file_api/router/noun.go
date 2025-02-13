package router

import (
	"LearningGuide/file_api/api"
	"github.com/OuterCyrex/Gorra/GorraAPI/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func NounRouter(R *gin.RouterGroup) {
	File := R.Group("/noun").Use(middlewares.TracerMiddleware(opentracing.GlobalTracer()))
	{
		File.POST("", api.NewNoun)
		File.GET("/list", api.NounList)
		File.GET("/desc/:id", api.GetNounDesc)
		File.DELETE("/:id", api.DeleteNoun)
	}
}
