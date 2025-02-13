package router

import (
	"LearningGuide/post_api/api"
	"github.com/OuterCyrex/Gorra/GorraAPI/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func FavRouter(R *gin.RouterGroup) {
	Like := R.Group("/fav").Use(middlewares.TracerMiddleware(opentracing.GlobalTracer()))
	{
		Like.POST("", api.NewFav)
		Like.GET("/list", api.GetPostByFav)
		Like.DELETE("", api.DeleteFav)
	}
}
