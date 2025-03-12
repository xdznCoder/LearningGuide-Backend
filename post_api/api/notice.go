package api

import (
	"LearningGuide/post_api/global"
	PostProto "LearningGuide/post_api/proto/.PostProto"
	UserProto "LearningGuide/post_api/proto/userProto"
	"errors"
	"github.com/OuterCyrex/Gorra/GorraAPI"
	err_resp "github.com/OuterCyrex/Gorra/GorraAPI/resp"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CheckNotice(c *gin.Context) {
	ctx := GorraAPI.RawContextWithSpan(c)

	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的路径参数",
		})
		return
	}

	checked, err := global.PostSrvClient.CheckNotice(ctx, &PostProto.CheckNoticeRequest{UserId: int32(userId)})

	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"new_notice": checked.NewNotices,
	})
}

func NoticeList(c *gin.Context) {
	userId, err1 := strconv.Atoi(c.DefaultQuery("user_id", "0"))
	noticeType, err2 := strconv.Atoi(c.DefaultQuery("type", "0"))
	pageNum, err3 := strconv.Atoi(c.DefaultQuery("pageNum", "0"))
	pageSize, err4 := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	err := errors.Join(err1, err2, err3, err4)
	if err != nil || userId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "无效的查询参数",
		})
		return
	}

	ctx := GorraAPI.RawContextWithSpan(c)

	notices, err := global.PostSrvClient.GetNoticeList(ctx, &PostProto.NoticeFilterRequest{
		UserId:   int32(userId),
		Type:     int32(noticeType),
		PageSize: int32(pageSize),
		PageNum:  int32(pageNum),
	})

	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	// 获取用户信息

	var ids []int32

	for _, n := range notices.Data {
		ids = append(ids, n.UserId)
	}

	userInfos, err := global.UserSrvClient.GetUsersByIds(ctx, &UserProto.IdsRequest{Ids: ids})
	if err != nil {
		err_resp.HandleGrpcErrorToHttp(err, c)
		return
	}

	userInfoMap := make(map[int32]*UserProto.UserInfoResponse)

	for _, u := range userInfos.Data {
		userInfoMap[u.Id] = u
	}

	// 将用户信息封装至消息列表中

	var respList []gin.H

	for _, n := range notices.Data {
		respList = append(respList, gin.H{
			"userInfo":   getUserDTO(userInfoMap[n.UserId]),
			"owner_id":   n.OwnerId,
			"type":       n.Type,
			"post_id":    n.PostId,
			"post_title": n.PostTitle,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  respList,
		"total": notices.Total,
	})
}

func getUserDTO(resp *UserProto.UserInfoResponse) map[string]any {
	return map[string]any{
		"id":       resp.Id,
		"email":    resp.Email,
		"nickname": resp.NickName,
		"gender":   resp.Gender,
		"birthday": resp.BirthDay,
		"image":    resp.Image,
		"desc":     resp.Desc,
	}
}
