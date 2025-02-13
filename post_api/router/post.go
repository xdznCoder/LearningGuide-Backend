package router

import (
	"LearningGuide/post_api/api"
	"LearningGuide/post_api/global"
	"github.com/OuterCyrex/Gorra/GorraAPI/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func PostRouter(R *gin.RouterGroup) {
	Post := R.Group("/post").Use(middlewares.TracerMiddleware(opentracing.GlobalTracer()))
	{
		Post.POST("", api.NewPost)
		Post.GET("/list", middlewares.JWTAuthMiddleware(global.ServerConfig.JwtKey), api.PostList)
		Post.GET("/:id", api.GetPost)
		Post.PUT("/:id", api.UpdatePost)
		Post.DELETE("/:id", api.DeletePost)
	}
}
