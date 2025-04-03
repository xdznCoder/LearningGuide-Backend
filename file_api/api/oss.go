package api

import (
	"LearningGuide/file_api/global"
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"io"
	"time"
)

type OssClientProxy struct {
	c *oss.Client
}

var OssClient *OssClientProxy

func (o *OssClientProxy) SetClient(c *oss.Client) {
	o.c = c
}

func (o *OssClientProxy) FileURL(ossName string, fileName string) (string, error) {
	expiration := time.Now().Add(1 * time.Hour)

	req := &oss.GetObjectRequest{
		Bucket: oss.Ptr(global.ServerConfig.AliyunOss.FileBucketName),
		Key:    oss.Ptr(ossName),
		RequestCommon: oss.RequestCommon{
			Parameters: map[string]string{
				"response-content-disposition": `attachment; filename="` + fileName + `"`,
			},
		},
	}

	signedURL, err := o.c.Presign(context.TODO(), req, oss.PresignExpiration(expiration))
	if err != nil {
		return "", err
	}

	return signedURL.URL, nil
}

func (o *OssClientProxy) Delete(ossName string) error {
	_, err := o.c.DeleteObject(context.Background(), &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(global.ServerConfig.AliyunOss.FileBucketName),
		Key:    oss.Ptr(ossName),
	})
	return err
}

func (o *OssClientProxy) Upload(ossName string, fileName string, file io.Reader) error {
	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr(global.ServerConfig.AliyunOss.FileBucketName),
		Key:    oss.Ptr(ossName),
		Body:   file,
		Metadata: map[string]string{
			"Content-Disposition": `attachment; filename="` + fileName + `"`,
		},
	}

	_, err := o.c.PutObject(context.TODO(), request)
	return err
}
