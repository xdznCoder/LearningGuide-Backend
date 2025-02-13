package router

import (
	"LearningGuide/post_api/api"
	"github.com/OuterCyrex/Gorra/GorraAPI/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func CommentRouter(R *gin.RouterGroup) {
	Comment := R.Group("/comment").Use(middlewares.TracerMiddleware(opentracing.GlobalTracer()))
	{
		Comment.POST("", api.NewComment)
		Comment.GET("/list", api.CommentList)
		Comment.DELETE("/:id", api.DeleteComment)
	}
}
