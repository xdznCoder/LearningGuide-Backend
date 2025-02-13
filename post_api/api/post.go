package api

import (
	"LearningGuide/post_api/forms"
	"LearningGuide/post_api/global"
	PostProto "LearningGuide/post_api/proto/.PostProto"
	"errors"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	err_resp "github.com/OuterCyrex/Gorra/GorraAPI/resp"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func PostList(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	uid, _ := c.Get("userId")

	userId, err := strconv.Atoi(c.DefaultQuery("user_id", "0"))
	title := c.DefaultQuery("title", "")
	category := c.DefaultQuery("category", "")
	pageNum, err := strconv.Atoi(c.DefaultQuery("pageNum", "0"))
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "0"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的查询参数",
		})
		return
	}

	resp, err := global.PostSrvClient.PostList(ctx, &PostProto.PostFilterRequest{
		UserId:   int32(userId),
		Title:    title,
		Category: category,
		PageNum:  int32(pageNum),
		PageSize: int32(pageSize),
	})

	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	var respList []gin.H

	for _, v := range resp.Data {
		liked, err := isLiked(&global.RDB, v.Id, int32(uid.(uint)))
		if !errors.Is(err, redis.Nil) && err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "服务器内部错误",
			})
			zap.S().Errorf("get post-like to redis failed: %v", err)
			return
		}

		favorite, err := isFav(&global.RDB, v.Id, int32(uid.(uint)))
		if !errors.Is(err, redis.Nil) && err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "服务器内部错误",
			})
			zap.S().Errorf("get post-fav to redis failed: %v", err)
			return
		}

		respList = append(respList, gin.H{
			"id":          v.Id,
			"title":       v.Title,
			"desc":        v.Desc,
			"user_id":     v.UserId,
			"image":       v.Image,
			"category":    v.Category,
			"comment_num": v.CommentNum,
			"fav_num":     v.FavNum,
			"like_num":    v.LikeNum,
			"is_liked":    liked,
			"is_favorite": favorite,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total": resp.Total,
		"data":  respList,
	})
}

func GetPost(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的路径参数",
		})
		return
	}

	resp, err := global.PostSrvClient.GetPost(ctx, &PostProto.PostID{Id: int32(id)})
	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func NewPost(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	var post forms.NewPostForm

	err := c.ShouldBindJSON(&post)
	if err != nil {
		err_resp.HandleValidatorError(err, c)
		return
	}

	resp, err := global.PostSrvClient.NewPost(ctx, &PostProto.NewPostRequest{
		UserId:   post.UserId,
		Category: post.Category,
		Content:  post.Content,
		Title:    post.Title,
		Desc:     post.Desc,
		Image:    post.Image,
	})

	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func UpdatePost(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的路径参数",
		})
		return
	}

	var post forms.UpdatePostForm

	err = c.ShouldBindJSON(&post)
	if err != nil {
		err_resp.HandleValidatorError(err, c)
		return
	}

	_, err = global.PostSrvClient.UpdatePost(ctx, &PostProto.UpdatePostRequest{
		Id:      int32(id),
		Content: post.Content,
		Title:   post.Title,
		Desc:    post.Desc,
		Image:   post.Image,
	})

	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

func DeletePost(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的路径参数",
		})
		return
	}

	_, err = global.PostSrvClient.DeletePost(ctx, &PostProto.DeletePostRequest{Id: int32(id)})
	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}
