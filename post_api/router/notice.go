package router

import (
	"LearningGuide/post_api/api"
	"github.com/OuterCyrex/Gorra/GorraAPI/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func NoticeRouter(R *gin.RouterGroup) {
	Notice := R.Group("/notice").Use(middlewares.TracerMiddleware(opentracing.GlobalTracer()))
	{
		Notice.GET("/:id", api.CheckNotice)
		Notice.GET("/list", api.NoticeList)
	}
}
