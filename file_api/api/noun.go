package api

import (
	"LearningGuide/file_api/forms"
	"LearningGuide/file_api/global"
	FileProto "LearningGuide/file_api/proto/.FileProto"
	"LearningGuide/file_api/utils"
	"context"
	"errors"
	"github.com/OuterCyrex/ChatGLM_sdk"
	GLM_Model "github.com/OuterCyrex/ChatGLM_sdk/model"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	handleGrpc "github.com/OuterCyrex/Gorra/GorraAPI/resp"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func NounList(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "0"))
	pageNum, err := strconv.Atoi(c.DefaultQuery("pageNum", "0"))
	courseId, err := strconv.Atoi(c.DefaultQuery("course_id", "0"))
	name := c.DefaultQuery("name", "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的查询参数",
		})
		return
	}

	resp, err := global.FileSrvClient.NounList(ctx, &FileProto.NounListRequest{
		Name:     name,
		CourseId: int32(courseId),
		PageNum:  int32(pageNum),
		PageSize: int32(pageSize),
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func DeleteNoun(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的路径参数",
		})
		return
	}

	_, err = global.FileSrvClient.DeleteNoun(ctx, &FileProto.DeleteNounRequest{Id: int32(id)})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func NewNoun(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	noun := forms.NewNounForm{}

	err := c.ShouldBindJSON(&noun)
	if err != nil {
		handleGrpc.HandleValidatorError(err, c)
		return
	}

	if len(noun.FileIds) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "文件最多可选10个",
		})
		return
	}

	var messageCtx *ChatGLM_sdk.MessageContext

	if len(noun.FileIds) == 0 {
		var ids []int32

		resp, err := global.FileSrvClient.FileList(ctx, &FileProto.FileFilterRequest{
			PageNum:  0,
			PageSize: 10,
			CourseId: noun.CourseId,
		})
		if err != nil {
			handleGrpc.HandleGrpcErrorToHttp(err, c)
			return
		}

		for _, v := range resp.Data {
			ids = append(ids, v.Id)
		}

		messageCtx, err = getFileContext(ctx, ids)
		if err != nil {
			handleGrpc.HandleGrpcErrorToHttp(err, c)
			return
		}
	} else {
		messageCtx, err = getFileContext(ctx, noun.FileIds)
		if err != nil {
			handleGrpc.HandleGrpcErrorToHttp(err, c)
			return
		}
	}

	client := ChatGLM_sdk.NewClient(global.ServerConfig.ChatGLM.AccessKey)

	id, err := client.SendAsync(messageCtx, getNounPrompt(noun.Name))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务器内部错误",
		})
		zap.S().Errorf("get async id from glm failed: %v", err)
		return
	}

	resp, err := global.FileSrvClient.NewNoun(ctx, &FileProto.NewNounRequest{
		Name:     noun.Name,
		Content:  id,
		CourseId: noun.CourseId,
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func GetNounDesc(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的路径参数",
		})
		return
	}

	resp, err := global.FileSrvClient.GetNounDetail(ctx, &FileProto.NounDetailRequest{Id: int32(id)})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	client := ChatGLM_sdk.NewClient(global.ServerConfig.ChatGLM.AccessKey)

	result := client.GetAsyncMessage(ChatGLM_sdk.NewContext(), resp.Content)

	if errors.Is(result.Error, ChatGLM_sdk.ErrResultProcessing) {
		c.JSON(http.StatusAccepted, gin.H{
			"msg": "GLM正在生成中",
		})
		return
	} else if result.Error != nil || len(result.Message) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务器内部错误",
		})
		zap.S().Errorf("get async message from glm failed: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content": result.Message[0].Content,
	})
}

func getFileContext(ctx context.Context, fileIds []int32) (*ChatGLM_sdk.MessageContext, error) {
	messageCtx := ChatGLM_sdk.NewContext()

	var files []*FileProto.FileInfoResponse

	for _, id := range fileIds {
		file, err := getFileInfo(ctx, int(id))
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	client := getOssClient(global.ServerConfig.AliyunOss)

	for _, file := range files {
		resp, err := client.GetObject(context.TODO(), &oss.GetObjectRequest{
			Bucket: oss.Ptr(global.ServerConfig.AliyunOss.BucketName),
			Key:    oss.Ptr(file.OssUrl),
		})

		if err != nil {
			return nil, err
		}

		content, err := utils.ReadFile(resp.Body, file.FileName)

		if err != nil {
			return nil, err
		}

		*messageCtx = append(*messageCtx, GLM_Model.Message{
			Role:    "user",
			Content: content,
		})
	}

	return messageCtx, nil
}

func getNounPrompt(name string) string {
	return "请根据提供的文件的内容来对" + name + "这一名词进行解释，要求对文件内容进行一定的隐式引用，且清晰简洁，一定要全面"
}
