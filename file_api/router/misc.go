package router

import (
	"LearningGuide/file_api/api"
	"github.com/gin-gonic/gin"
)

func MiscRouter(R *gin.RouterGroup) {
	Misc := R.Group("")
	{
		Misc.POST("/picture", api.UploadPictures)
	}
}
