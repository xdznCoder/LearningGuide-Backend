package api

import (
	"LearningGuide/file_api/forms"
	"LearningGuide/file_api/global"
	ChatProto "LearningGuide/file_api/proto/.ChatProto"
	FileProto "LearningGuide/file_api/proto/.FileProto"
	"fmt"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	handleGrpc "github.com/OuterCyrex/Gorra/GorraAPI/resp"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"strconv"
	"strings"
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

func SendMessage(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	message := forms.SendMessageForm{}

	err := c.ShouldBindJSON(&message)

	if err != nil {
		handleGrpc.HandleValidatorError(err, c)
		return
	}

	var stream grpc.ServerStreamingClient[ChatProto.ChatModelResponse]
	query := "你好"

	switch message.Type {
	case 1:
		stream, err = global.ChatSrvClient.SendStreamMessage(ctx, &ChatProto.UserMessage{
			CourseID:     message.CourseId,
			SessionID:    message.SessionId,
			Content:      message.Content,
			FileURL:      "",
			TemplateType: int32(TemplateTypeUserQuery),
		})
		if err != nil {
			handleGrpc.HandleGrpcErrorToHttp(err, c)
			return
		}
		query = message.Content
	case 2:
		iErr := error(nil)
		fileId, iErr := strconv.Atoi(message.Content)
		if iErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "无效文件ID",
			})
			return
		}

		resp, iErr := getFileInfo(ctx, fileId)

		if iErr != nil {
			handleGrpc.HandleGrpcErrorToHttp(err, c)
			return
		}

		url, iErr := OssClient.FileURL(resp.OssUrl, resp.FileName)
		if iErr != nil {
			handleGrpc.HandleGrpcErrorToHttp(err, c)
			return
		}

		if _, err = global.FileSrvClient.NewMessage(ctx, &FileProto.NewMessageRequest{
			Content:   url,
			Type:      int32(message.Type),
			SessionId: message.SessionId,
			Speaker:   "user",
		}); err != nil {
			handleGrpc.HandleGrpcErrorToHttp(err, c)
			return
		}

		if _, err = global.FileSrvClient.NewMessage(ctx, &FileProto.NewMessageRequest{
			Content:   "我已收到文件，请您询问关于文件的任何问题",
			Type:      1,
			SessionId: message.SessionId,
			Speaker:   "assistant",
		}); err != nil {
			handleGrpc.HandleGrpcErrorToHttp(err, c)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "文件上传成功",
		})
		return
	}

	if stream == nil {
		return
	}

	var result strings.Builder

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	for {
		output, iErr := stream.Recv()
		if iErr == io.EOF {
			break
		} else if iErr != nil {
			handleGrpc.HandleGrpcErrorToHttp(err, c)
			return
		}

		_, _ = fmt.Fprintf(c.Writer, "data: %s\n", output.Content)
		result.WriteString(output.Content)
		if flusher, ok := c.Writer.(http.Flusher); ok {
			flusher.Flush()
		}
		time.Sleep(10 * time.Millisecond)
	}

	if _, err = global.FileSrvClient.NewMessage(ctx, &FileProto.NewMessageRequest{
		Content:   query,
		Type:      int32(message.Type),
		SessionId: message.SessionId,
		Speaker:   "user",
	}); err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	_, err = global.FileSrvClient.NewMessage(ctx, &FileProto.NewMessageRequest{
		Content:   result.String(),
		Type:      int32(message.Type),
		SessionId: message.SessionId,
		Speaker:   "assistant",
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}
}
