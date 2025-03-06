package logic

import (
	"LearningGuide/user_srv/model"
	"LearningGuide/user_srv/userProto"
)

func ModelToResponse(user model.User) *userProto.UserInfoResponse {
	var birthday uint64
	if user.Birthday != nil {
		birthday = uint64(user.Birthday.Unix())
	}
	return &userProto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		Email:    user.Email,
		NickName: user.NickName,
		BirthDay: birthday,
		Gender:   user.Gender,
		Role:     int32(user.Role),
		Image:    user.Image,
		Desc:     user.Desc,
	}
}
