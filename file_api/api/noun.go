package api

import (
	"LearningGuide/file_api/forms"
	"LearningGuide/file_api/global"
	ChatProto "LearningGuide/file_api/proto/.ChatProto"
	FileProto "LearningGuide/file_api/proto/.FileProto"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	handleGrpc "github.com/OuterCyrex/Gorra/GorraAPI/resp"
	"github.com/gin-gonic/gin"
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

	stream, err := global.ChatSrvClient.SendStreamMessage(ctx, &ChatProto.UserMessage{
		CourseID:     noun.CourseId,
		Content:      noun.Name,
		TemplateType: int32(TemplateTypeNounExplainGenerate),
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	result, err := ToString(stream)

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	resp, err := global.FileSrvClient.NewNoun(ctx, &FileProto.NewNounRequest{
		Name:     noun.Name,
		Content:  result,
		CourseId: noun.CourseId,
	})

	if err != nil {
		handleGrpc.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      resp.Id,
		"content": result,
	})
}

func GetNounDetail(c *gin.Context) {
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

	c.JSON(http.StatusOK, resp)
}
