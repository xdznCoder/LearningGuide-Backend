package api

import (
	"LearningGuide/file_api/global"
	ChatProto "LearningGuide/file_api/proto/.ChatProto"
	FileProto "LearningGuide/file_api/proto/.FileProto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	handleGrpc "github.com/OuterCyrex/Gorra/GorraAPI/resp"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func FileList(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	fileName := c.DefaultQuery("file_name", "")
	fileType := c.DefaultQuery("file_type", "")
	userId, err1 := strconv.Atoi(c.DefaultQuery("user_id", "0"))
	courseId, err2 := strconv.Atoi(c.DefaultQuery("course_id", "0"))
	pageNum, err3 := strconv.Atoi(c.DefaultQuery("pageNum", "0"))
	pageSize, err4 := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	err := errors.Join(err1, err2, err3, err4)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效查询参数",
		})
		return
	}

	resp, err := global.FileSrvClient.FileList(ctx, &FileProto.FileFilterRequest{
		FileName: fileName,
		FileType: fileType,
		UserId:   int32(userId),
		CourseId: int32(courseId),
		PageNum:  int32(pageNum),
		PageSize: int32(pageSize),
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": resp.Total,
		"data":  resp.Data,
	})
}

func UploadFile(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的文件类型",
		})
		return
	}

	userId, err := strconv.Atoi(c.DefaultPostForm("user_id", "0"))
	if err != nil || userId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的ID参数",
		})
		return
	}

	courseId, err := strconv.Atoi(c.DefaultPostForm("course_id", "0"))

	if err != nil || courseId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的ID参数",
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

	ossName := fmt.Sprintf("%d-%d-%s", userId, courseId, fileHeader.Filename)

	err = OssClient.Upload(ossName, fileHeader.Filename, file)

	if err != nil {
		zap.S().Errorf("文件上传失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "文件上传失败",
		})
		return
	}

	resp, err := global.FileSrvClient.CreateFile(ctx, &FileProto.CreateFileRequest{
		FileName: fileHeader.Filename,
		FileType: filepath.Ext(fileHeader.Filename),
		FileSize: fileHeader.Size,
		OssUrl:   ossName,
		UserId:   int32(userId),
		CourseId: int32(courseId),
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	global.RDB.Del(context.TODO(), fmt.Sprintf("%d", resp.Id))
	_, err = getFileInfo(ctx, int(resp.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "文件上传失败",
		})
		zap.S().Errorf("set fileInfo to redis failed: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": resp.Id,
	})
}

func GetFileDetail(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	id := c.Param("id")

	fileId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效路径参数",
		})
		return
	}

	resp, err := getFileInfo(ctx, fileId)

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func DownloadFile(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	id := c.Param("id")

	fileId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效路径参数",
		})
		return
	}

	resp, err := getFileInfo(ctx, fileId)

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	url, err := OssClient.FileURL(resp.OssUrl, resp.FileName)
	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

func UpdateFileDesc(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	id := c.Param("id")

	fileId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效路径参数",
		})
		return
	}

	resp, err := getFileInfo(ctx, fileId)

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	url, err := OssClient.FileURL(resp.OssUrl, resp.FileName)

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	stream, err := global.ChatSrvClient.SendStreamMessage(ctx, &ChatProto.UserMessage{
		CourseID:     resp.CourseId,
		Content:      "",
		FileURL:      url,
		TemplateType: int32(TemplateTypeFileDescribeGenerate),
	})

	result, err := ToString(stream)

	_, err = global.FileSrvClient.UpdateFile(ctx, &FileProto.UpdateFileRequest{
		Id:   int32(fileId),
		Desc: result,
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	global.RDB.Del(context.Background(), fmt.Sprintf("%d", fileId))

	c.JSON(http.StatusOK, gin.H{
		"Desc": result,
	})
}

func DeleteFile(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	id := c.Param("id")

	fileId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效路径参数",
		})
		return
	}

	resp, err := getFileInfo(ctx, fileId)

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	_, err = global.FileSrvClient.DeleteFile(ctx, &FileProto.DeleteFileRequest{Id: int32(fileId)})
	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	global.RDB.Del(ctx, fmt.Sprintf("%d", fileId))

	err = OssClient.Delete(resp.OssUrl)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务器内部错误",
		})
		zap.S().Errorf("oss delete object failed: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func UpdateFileMindMap(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	fileId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的路径参数",
		})
		return
	}

	resp, err := getFileInfo(ctx, fileId)

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	url, err := OssClient.FileURL(resp.OssUrl, resp.FileName)

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	stream, err := global.ChatSrvClient.SendStreamMessage(ctx, &ChatProto.UserMessage{
		CourseID:     resp.CourseId,
		Content:      "",
		FileURL:      url,
		TemplateType: int32(TemplateTypeMindMapGenerate),
	})

	result, err := ToString(stream)

	mindMap := transResultToStringJSON(result)

	_, err = global.FileSrvClient.UpdateFile(ctx, &FileProto.UpdateFileRequest{
		Id:      int32(fileId),
		MindMap: mindMap,
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	global.RDB.Del(context.Background(), fmt.Sprintf("%d", fileId))

	c.JSON(http.StatusOK, gin.H{
		"mind_map": mindMap,
	})
}

func getFileInfo(ctx context.Context, id int) (*FileProto.FileInfoResponse, error) {
	result, err := global.RDB.Get(ctx, fmt.Sprintf("%d", id)).Result()

	if errors.Is(err, redis.Nil) {
		resp, rpcErr := global.FileSrvClient.GetFileDetail(ctx, &FileProto.FileDetailRequest{Id: int32(id)})
		if rpcErr != nil {
			return nil, rpcErr
		}

		fileInfo, err := json.Marshal(resp)
		if err != nil {
			return nil, err
		}

		err = global.RDB.Set(ctx, fmt.Sprintf("%d", id), fileInfo, 20*time.Minute).Err()
		if err != nil {
			return nil, xerrors.Errorf("failed to set file name in Redis: %v", err)
		}

		return resp, nil
	} else if err != nil {
		return nil, xerrors.Errorf("failed to get file name in Redis: %v", err)
	}

	var fileInfo FileProto.FileInfoResponse

	err = json.Unmarshal([]byte(result), &fileInfo)
	if err != nil {
		return nil, err
	}

	return &fileInfo, nil
}
