package api

import (
	"LearningGuide/file_api/forms"
	"LearningGuide/file_api/global"
	FileProto "LearningGuide/file_api/proto/.FileProto"
	"LearningGuide/file_api/utils"
	"context"
	"fmt"
	"github.com/OuterCyrex/ChatGLM_sdk"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	handleGrpc "github.com/OuterCyrex/Gorra/GorraAPI/resp"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func SessionList(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	courseId, err := strconv.Atoi(c.DefaultQuery("course_id", "0"))
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "0"))
	pageNum, err := strconv.Atoi(c.DefaultQuery("pageNum", "0"))

	resp, err := global.FileSrvClient.SessionList(ctx, &FileProto.SessionListRequest{
		CourseId: int32(courseId),
		PageSize: int32(pageSize),
		PageNum:  int32(pageNum),
	})
	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func NewSession(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	var SessionForm forms.CreateSessionForm

	err := c.ShouldBindJSON(&SessionForm)
	if err != nil {
		handleGrpc.HandleValidatorError(err, c)
		return
	}

	resp, err := global.FileSrvClient.CreateSession(ctx, &FileProto.CreateSessionRequest{
		CourseId: SessionForm.CourseId,
	})
	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	sessionManager.New(resp.Id)

	c.JSON(http.StatusOK, resp)
}

func DeleteSession(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	Id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的路径参数",
		})
		return
	}

	sessionManager.Del(int32(Id))

	_, err = global.FileSrvClient.DeleteSession(ctx, &FileProto.DeleteSessionRequest{Id: int32(Id)})
	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func MessageList(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	sessionId, err := strconv.Atoi(c.DefaultQuery("session_id", "0"))
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "0"))
	pageNum, err := strconv.Atoi(c.DefaultQuery("pageNum", "0"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效查询参数",
		})
		return
	}

	resp, err := global.FileSrvClient.MessageList(ctx, &FileProto.MessageListRequest{
		SessionId: int32(sessionId),
		PageSize:  int32(pageSize),
		PageNum:   int32(pageNum),
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

type SessionManager struct {
	lock    sync.Mutex
	session map[int32]*ChatGLM_sdk.MessageContext
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		lock:    sync.Mutex{},
		session: make(map[int32]*ChatGLM_sdk.MessageContext),
	}
}

func (s *SessionManager) New(Id int32) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.session[Id] = ChatGLM_sdk.NewContext()
}

func (s *SessionManager) Del(Id int32) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.session, Id)
}

var sessionManager *SessionManager

func init() {
	sessionManager = NewSessionManager()
}

func SendMessage(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	message := forms.SendMessageForm{}

	err := c.ShouldBindJSON(&message)

	if err != nil {
		handleGrpc.HandleValidatorError(err, c)
		return
	}

	// 参数校验
	if _, ok := sessionManager.session[message.SessionId]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "会话次数已达上限，请开始新的会话",
		})
		return
	}

	if len(*sessionManager.session[message.SessionId]) > 30 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "会话次数已达上限，请开始新的会话",
		})

		sessionManager.Del(message.SessionId)

		return
	}

	// 模式选择

	client := ChatGLM_sdk.NewClient(global.ServerConfig.ChatGLM.AccessKey)

	msgChannel := make(<-chan ChatGLM_sdk.Result)

	var resp strings.Builder

	switch message.Type {
	case 1:
		msgChannel = client.SendStream(sessionManager.session[message.SessionId], message.Content)
	case 2:
		fileId, err := strconv.Atoi(message.Content)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "无效文件ID",
			})
			return
		}

		resp, err := getFileInfo(ctx, fileId)

		if err != nil {
			handleGrpc.HandleGrpcErrorToHttp(err, c)
			return
		}

		ossClient := getOssClient(global.ServerConfig.AliyunOss)

		result, err := ossClient.GetObject(context.TODO(), &oss.GetObjectRequest{
			Bucket: oss.Ptr(global.ServerConfig.AliyunOss.FileBucketName),
			Key:    oss.Ptr(resp.OssUrl),
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "服务器内部错误",
			})
			zap.S().Errorf("get object from oss failed: %v", err)
			return
		}

		file, err := utils.ReadFile(result.Body, resp.FileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "服务器内部错误",
			})
			zap.S().Errorf("read from file failed: %v", err)
			return
		}

		msgChannel = client.SendStream(sessionManager.session[message.SessionId], file)
	}

	for words := range msgChannel {
		if len(words.Message) > 0 {

			content := words.Message[0].Content

			resp.WriteString(content)

			_, _ = fmt.Fprintf(c.Writer, "data: %s\n", content)
			if flusher, ok := c.Writer.(http.Flusher); ok {
				flusher.Flush()
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	if _, err = global.FileSrvClient.NewMessage(ctx, &FileProto.NewMessageRequest{
		Content:   message.Content,
		Type:      int32(message.Type),
		SessionId: message.SessionId,
		Speaker:   "user",
	}); err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	_, err = global.FileSrvClient.NewMessage(ctx, &FileProto.NewMessageRequest{
		Content:   resp.String(),
		Type:      int32(message.Type),
		SessionId: message.SessionId,
		Speaker:   "assistant",
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}
}
