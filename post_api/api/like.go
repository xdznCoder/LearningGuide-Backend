package api

import (
	"LearningGuide/post_api/forms"
	"LearningGuide/post_api/global"
	PostProto "LearningGuide/post_api/proto/.PostProto"
	"context"
	"errors"
	"fmt"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	err_resp "github.com/OuterCyrex/Gorra/GorraAPI/resp"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

func NewLike(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	var like forms.NewLikeForm
	err := c.ShouldBindJSON(&like)

	if err != nil {
		err_resp.HandleValidatorError(err, c)
		return
	}

	liked, err := isLiked(&global.RDB, like.PostId, like.UserId)
	if liked && err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "用户已点赞",
		})
		return
	} else if err != nil && !errors.Is(err, redis.Nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务器内部错误",
		})
		zap.S().Errorf("get post-like to redis failed: %v", err)
		return
	}

	_, err = global.PostSrvClient.NewLike(ctx, &PostProto.NewLikeRequest{
		PostId: like.PostId,
		UserId: like.UserId,
	})

	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	err = likePost(&global.RDB, like.PostId, like.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务器内部错误",
		})
		zap.S().Errorf("add post-like to redis failed: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "点赞成功",
	})
}

func GetPostByLike(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	userId, err := strconv.Atoi(c.DefaultQuery("user_id", "0"))
	pageNum, err := strconv.Atoi(c.DefaultQuery("pageNum", "0"))
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "0"))

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "无效的查询参数",
		})
		return
	}

	resp, err := global.PostSrvClient.PostListByLike(ctx, &PostProto.LikeListRequest{
		UserId:   int32(userId),
		PageNum:  int32(pageNum),
		PageSize: int32(pageSize),
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

func DeleteLike(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	var like forms.DeleteLikeForm

	err := c.ShouldBindJSON(&like)

	if err != nil {
		err_resp.HandleValidatorError(err, c)
		return
	}

	err = cancelLike(&global.RDB, like.PostId, like.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务器内部错误",
		})
		zap.S().Errorf("delete post-like to redis failed: %v", err)
		return
	}

	_, err = global.PostSrvClient.DeleteLike(ctx, &PostProto.DeleteLikeRequest{
		UserId: like.UserId,
		PostId: like.PostId,
	})

	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func likePost(rdb *redis.Client, postId int32, userId int32) error {
	key := fmt.Sprintf("post-likes:%d", postId)
	score := float64(time.Now().UnixNano() / int64(time.Millisecond))
	_, err := rdb.ZAdd(context.Background(), key, &redis.Z{Score: score, Member: userId}).Result()
	return err
}

func cancelLike(rdb *redis.Client, postId int32, userId int32) error {
	key := fmt.Sprintf("post-likes:%d", postId)
	_, err := rdb.ZRem(context.Background(), key, userId).Result()
	return err
}

func isLiked(rdb *redis.Client, postId int32, userId int32) (bool, error) {
	key := fmt.Sprintf("post-likes:%d", postId)
	score, err := rdb.ZScore(context.Background(), key, fmt.Sprint(userId)).Result()

	if err != nil {
		return false, err
	}

	if score != 0 {
		return true, nil
	}

	return false, nil
}
