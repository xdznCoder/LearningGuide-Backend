package initialize

import (
	"LearningGuide/user_api/middlewares"
	"LearningGuide/user_api/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Routers() *gin.Engine {
	R := gin.Default()
	R.Use(middlewares.Cors())

	// keep-alive 检查API是否存活
	R.GET("health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
		})
	})

	ApiGroup := R.Group("/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)

	return R
}
