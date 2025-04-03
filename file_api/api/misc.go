package api

import (
	"LearningGuide/file_api/global"
	"context"
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func UploadPictures(c *gin.Context) {
	fileHeader, err := c.FormFile("picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的图片类型",
		})
		return
	}

	if fileHeader.Size >= 5242880 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "文件大小超过5MB",
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的文件类型",
		})
		return
	}

	cfg := oss.LoadDefaultConfig().WithCredentialsProvider(
		credentials.NewStaticCredentialsProvider(
			global.ServerConfig.AliyunOss.AccessKey,
			global.ServerConfig.AliyunOss.SecretKey,
			"")).
		WithRegion(global.ServerConfig.AliyunOss.Region)

	client := oss.NewClient(cfg)

	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr(global.ServerConfig.AliyunOss.PictureBucketName),
		Key:    oss.Ptr(fileHeader.Filename),
		Body:   file,
	}

	_, err = client.PutObject(context.TODO(), request)

	if err != nil {
		zap.S().Errorf("文件上传失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "文件上传失败",
		})
		return
	}

	url := fmt.Sprintf("https://%s.%s/%s",
		global.ServerConfig.AliyunOss.PictureBucketName,
		global.ServerConfig.AliyunOss.EndPoint,
		fileHeader.Filename,
	)

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}
