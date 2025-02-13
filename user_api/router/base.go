package router

import (
	"LearningGuide/user_api/api"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	zap.S().Infof("Initialize LearningGuide User Router...")
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("/captcha", api.GetCaptcha)
		BaseRouter.POST("/sendEmail", api.SendEmail)
	}
}
