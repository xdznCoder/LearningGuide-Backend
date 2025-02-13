package router

import (
	"LearningGuide/file_api/api"
	"github.com/OuterCyrex/Gorra/GorraAPI/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func MessageRouter(R *gin.RouterGroup) {
	Session := R.Group("/session").Use(middlewares.TracerMiddleware(opentracing.GlobalTracer()))
	{
		Session.GET("/list", api.SessionList)
		Session.POST("", api.NewSession)
		Session.DELETE("/:id", api.DeleteSession)
	}

	Message := R.Group("/message").Use(middlewares.TracerMiddleware(opentracing.GlobalTracer()))
	{
		Message.GET("/list", api.MessageList)
		Message.POST("/send", api.SendMessage)
		Message.GET("/new", api.SetUpWebsocket)
	}
}
