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

func NewFav(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	var fav forms.NewFavForm
	err := c.ShouldBindJSON(&fav)

	if err != nil {
		err_resp.HandleValidatorError(err, c)
		return
	}

	favorite, err := isFav(&global.RDB, fav.PostId, fav.UserId)
	if favorite && err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "用户已收藏",
		})
		return
	} else if err != nil && !errors.Is(err, redis.Nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务器内部错误",
		})
		zap.S().Errorf("get post-fav to redis failed: %v", err)
		return
	}

	_, err = global.PostSrvClient.NewFav(ctx, &PostProto.NewFavRequest{
		PostId: fav.PostId,
		UserId: fav.UserId,
	})

	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	err = favPost(&global.RDB, fav.PostId, fav.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务器内部错误",
		})
		zap.S().Errorf("add post-like to redis failed: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "收藏成功",
	})
}

func GetPostByFav(c *gin.Context) {
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

	resp, err := global.PostSrvClient.PostListByFav(ctx, &PostProto.FavListRequest{
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

func DeleteFav(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	var fav forms.DeleteFavForm

	err := c.ShouldBindJSON(&fav)

	if err != nil {
		err_resp.HandleValidatorError(err, c)
		return
	}

	err = cancelFav(&global.RDB, fav.PostId, fav.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "服务器内部错误",
		})
		zap.S().Errorf("delete post-fav to redis failed: %v", err)
		return
	}

	_, err = global.PostSrvClient.DeleteFav(ctx, &PostProto.DeleteFavRequest{
		UserId: fav.UserId,
		PostId: fav.PostId,
	})

	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func favPost(rdb *redis.Client, postId int32, userId int32) error {
	key := fmt.Sprintf("post-favs:%d", postId)
	score := float64(time.Now().UnixNano() / int64(time.Millisecond))
	_, err := rdb.ZAdd(context.Background(), key, &redis.Z{Score: score, Member: userId}).Result()
	return err
}

func cancelFav(rdb *redis.Client, postId int32, userId int32) error {
	key := fmt.Sprintf("post-favs:%d", postId)
	_, err := rdb.ZRem(context.Background(), key, userId).Result()
	return err
}

func isFav(rdb *redis.Client, postId int32, userId int32) (bool, error) {
	key := fmt.Sprintf("post-favs:%d", postId)
	score, err := rdb.ZScore(context.Background(), key, fmt.Sprint(userId)).Result()
	if err != nil {
		return false, err
	}

	if score != 0 {
		return true, nil
	}

	return false, nil
}
