package router

import (
	"LearningGuide/post_api/api"
	"github.com/OuterCyrex/Gorra/GorraAPI/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func LikeRouter(R *gin.RouterGroup) {
	Like := R.Group("/like").Use(middlewares.TracerMiddleware(opentracing.GlobalTracer()))
	{
		Like.POST("", api.NewLike)
		Like.GET("/list", api.GetPostByLike)
		Like.DELETE("", api.DeleteLike)
	}
}
