package api

import (
	"LearningGuide/post_api/global"
	UserProto "LearningGuide/post_api/proto/.UserProto"
	"context"
	"fmt"
)

func checkUserId(id int) bool {
	ctx := context.Background()

	isMember, err := global.RDB.SIsMember(ctx, "userIDSet", fmt.Sprintf("%d", id)).Result()
	if isMember {
		return true
	}

	_, err = global.UserSrvClient.GetUserById(ctx, &UserProto.IdRequest{Id: int32(id)})
	if err != nil {
		return false
	} else {
		_, _ = global.RDB.SAdd(ctx, "userIDSet", fmt.Sprintf("%d", id)).Result()
		return true
	}
}
