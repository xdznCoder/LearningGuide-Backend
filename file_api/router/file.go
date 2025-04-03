package router

import (
	"LearningGuide/file_api/api"
	"github.com/OuterCyrex/Gorra/GorraAPI/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func FileRouter(R *gin.RouterGroup) {
	File := R.Group("/file").Use(middlewares.TracerMiddleware(opentracing.GlobalTracer()))
	{
		File.POST("/upload", api.UploadFile)
		File.GET("/list", api.FileList)
		File.GET("/detail/:id", api.GetFileDetail)
		File.GET("/download/:id", api.DownloadFile)
		File.PUT("/desc/:id", api.UpdateFileDesc)
		File.DELETE("/:id", api.DeleteFile)
		File.PUT("/map/:id", api.UpdateFileMindMap)
	}
}
