package api

import (
	"LearningGuide/file_api/forms"
	"LearningGuide/file_api/global"
	FileProto "LearningGuide/file_api/proto/.FileProto"
	"LearningGuide/file_api/utils"
	"context"
	"errors"
	"github.com/OuterCyrex/ChatGLM_sdk"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	handleGrpc "github.com/OuterCyrex/Gorra/GorraAPI/resp"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"sync"
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

	SessionMapMu.Lock()
	delete(SessionMap, int32(Id))
	SessionMapMu.Unlock()

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

// websocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type session struct {
	history  *ChatGLM_sdk.MessageContext
	wsClient *websocket.Conn
}

var SessionMapMu sync.Mutex

var SessionMap = make(map[int32]session)

func SetUpWebsocket(c *gin.Context) {
	sessionId, err := strconv.Atoi(c.DefaultQuery("session_id", "0"))

	if err != nil || sessionId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的会话ID",
		})
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.S().Errorf("upgrade to WebSocket failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "建立连接失败",
		})
		return
	}

	SessionMapMu.Lock()
	SessionMap[int32(sessionId)] = session{
		wsClient: ws,
		history:  ChatGLM_sdk.NewContext(),
	}
	SessionMapMu.Unlock()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			zap.S().Infof("WebSocket connection closed for session %d", sessionId)
			return
		}
	}
}

func SendMessage(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	message := forms.SendMessageForm{}

	err := c.ShouldBindJSON(&message)

	if err != nil {
		handleGrpc.HandleValidatorError(err, c)
		return
	}

	if _, ok := SessionMap[message.SessionId]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "会话次数已达上限，请开始新的会话",
		})
		return
	}

	if len(*SessionMap[message.SessionId].history) > 30 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "会话次数已达上限，请开始新的会话",
		})

		SessionMapMu.Lock()
		delete(SessionMap, message.SessionId)
		SessionMapMu.Unlock()

		return
	}

	ws := SessionMap[message.SessionId].wsClient

	client := ChatGLM_sdk.NewClient(global.ServerConfig.ChatGLM.AccessKey)

	msgChannel := make(<-chan ChatGLM_sdk.Result)

	var resp strings.Builder

	switch message.Type {
	case 1:
		msgChannel = client.SendStream(SessionMap[message.SessionId].history, message.Content)
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

		msgChannel = client.SendStream(SessionMap[message.SessionId].history, file)
	}

	for words := range msgChannel {
		if len(words.Message) > 0 {
			resp.WriteString(words.Message[0].Content)
			err = ws.WriteMessage(websocket.TextMessage, []byte(words.Message[0].Content))
			if errors.Is(err, websocket.ErrCloseSent) {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "先建立websocket连接",
				})
				return
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "服务器内部错误",
				})
				zap.S().Errorf("write message to file failed: %v", err)
				return
			}
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

	c.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
