package initialize

import (
	"LearningGuide/file_api/api"
	"LearningGuide/file_api/config"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

func InitClient(c *config.OssConfig) {
	var client *oss.Client
	cfg := oss.LoadDefaultConfig().WithCredentialsProvider(
		credentials.NewStaticCredentialsProvider(c.AccessKey, c.SecretKey, "")).
		WithRegion(c.Region)

	client = oss.NewClient(cfg)

	api.OssClient = &api.OssClientProxy{}
	api.OssClient.SetClient(client)
}
