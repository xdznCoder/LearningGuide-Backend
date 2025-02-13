package api

import (
	"LearningGuide/post_api/forms"
	"LearningGuide/post_api/global"
	PostProto "LearningGuide/post_api/proto/.PostProto"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	err_resp "github.com/OuterCyrex/Gorra/GorraAPI/resp"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CommentList(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	postId, err := strconv.Atoi(c.DefaultQuery("post_id", "0"))
	userId, err := strconv.Atoi(c.DefaultQuery("user_id", "0"))
	parentId, err := strconv.Atoi(c.DefaultQuery("parent_id", "0"))
	pageNum, err := strconv.Atoi(c.DefaultQuery("pageNum", "0"))
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "0"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的查询参数",
		})
		return
	}

	resp, err := global.PostSrvClient.CommentList(ctx, &PostProto.CommentFilterRequest{
		UserId:          int32(userId),
		PostId:          int32(postId),
		ParentCommendId: int32(parentId),
		PageNum:         int32(pageNum),
		PageSize:        int32(pageSize),
	})

	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": resp.Total,
		"data":  resp.Data,
	})
}

func NewComment(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	var comment forms.NewCommentForm

	err := c.ShouldBindJSON(&comment)
	if err != nil {
		err_resp.HandleValidatorError(err, c)
		return
	}

	resp, err := global.PostSrvClient.NewComment(ctx, &PostProto.NewCommentRequest{
		UserId:          comment.UserId,
		PostId:          comment.PostId,
		ParentCommentId: comment.ParentCommentId,
		Content:         comment.Content,
	})
	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func DeleteComment(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的路径参数",
		})
		return
	}

	_, err = global.PostSrvClient.DeleteComment(ctx, &PostProto.DeleteCommentRequest{Id: int32(id)})
	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}
