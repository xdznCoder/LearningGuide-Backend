package utils

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				zap.S().Errorf("服务器内部错误: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "服务器内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": e.Message(),
				})
			case codes.Unavailable:
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"msg": "连接至rpc服务器错误",
				})
			case codes.AlreadyExists:
				c.JSON(http.StatusConflict, gin.H{
					"msg": e.Message(),
				})
			default:
				zap.S().Errorf("服务器内部错误: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Message(),
				})
			}
		} else {
			zap.S().Errorf("服务器内部错误: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "服务器错误",
			})
		}
	}
}
