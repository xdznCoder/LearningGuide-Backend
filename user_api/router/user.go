package router

import (
	"LearningGuide/user_api/api"
	"LearningGuide/user_api/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func InitUserRouter(Router *gin.RouterGroup) {
	zap.S().Infof("Initialize LearningGuide User Router...")
	UserRouter := Router.Group("/user").Use(middlewares.TracerMiddleware(opentracing.GlobalTracer()))
	{
		UserRouter.GET("list", middlewares.JWTAuthMiddleware(), api.GetUserList)
		UserRouter.POST("/register", api.Register)
		UserRouter.PUT("/:id", middlewares.JWTAuthMiddleware(), api.UpdateUser)
		UserRouter.DELETE("/:id", middlewares.JWTAuthMiddleware(), api.DeleteUser)
		UserRouter.POST("/login", api.PasswordLogin)
		UserRouter.PUT("/password/:id", middlewares.JWTAuthMiddleware(), api.ChangePassword)
		UserRouter.GET("/:id", middlewares.JWTAuthMiddleware(), api.GetUserDetail)
		UserRouter.GET("", middlewares.JWTAuthMiddleware(), api.GetUserByToken)
	}
}
