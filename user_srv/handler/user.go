package handler

import (
	"LearningGuide/user_srv/global"
	"LearningGuide/user_srv/model"
	proto "LearningGuide/user_srv/proto/.UserProto"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func ModelToResponse(user model.User) *proto.UserInfoResponse {
	var birthday uint64
	if user.Birthday != nil {
		birthday = uint64(user.Birthday.Unix())
	}
	return &proto.UserInfoResponse{
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

func Paginate(pageNum, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNum == 0 {
			pageNum = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (u UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	userListResponse := make([]*proto.UserInfoResponse, 0)

	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	err := global.DB.Scopes(Paginate(int(req.PageNum), int(req.PageSize))).Find(&users).Error
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		userListResponse = append(userListResponse, ModelToResponse(user))
	}

	return &proto.UserListResponse{
		Total: int32(result.RowsAffected),
		Data:  userListResponse,
	}, nil
}

func (u UserServer) GetUserByEmail(ctx context.Context, req *proto.EmailRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Email: req.Email}).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效用户ID")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return ModelToResponse(user), nil
}

func (u UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{BaseModel: model.BaseModel{ID: req.Id}}).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效用户ID")
	}
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return ModelToResponse(user), nil
}

func (u UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var formerUser model.User
	result := global.DB.Where(&model.User{Email: req.Email}).Find(&formerUser)
	if result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "电子邮箱已被使用")
	}

	salt, encodedPwd := password.Encode(req.GetPassword(), &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha256.New,
	})

	user := model.User{
		Email:    req.Email,
		Password: fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd),
		NickName: req.GetNickName(),
	}

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建用户失败: %v", result.Error)
	}

	return ModelToResponse(user), nil
}

func (u UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*proto.Empty, error) {
	var user model.User
	result := global.DB.Where(&model.User{BaseModel: model.BaseModel{ID: req.Id}}).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效用户ID")
	}
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	birthDay := time.Unix(int64(req.GetBirthDay()), 0)
	user.NickName = req.GetNickName()
	user.Birthday = &birthDay
	user.Gender = req.GetGender()
	user.Desc = req.GetDesc()
	user.Image = req.GetImage()

	result = global.DB.Updates(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "更新用户失败: %v", result.Error)
	}

	return &proto.Empty{}, nil
}

func (u UserServer) CheckPasswordInfo(ctx context.Context, req *proto.PasswordCheck) (*proto.CheckResponse, error) {
	passwordInfo := strings.Split(req.GetEncryptedPassword(), "$")
	check := password.Verify(req.GetPassword(), passwordInfo[2], passwordInfo[3], &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha256.New,
	})

	return &proto.CheckResponse{Success: check}, nil
}

func (u UserServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.Empty, error) {
	result := global.DB.Where(&model.User{BaseModel: model.BaseModel{ID: req.Id}}).Delete(&model.User{})
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效用户ID")
	}

	return &proto.Empty{}, nil
}

func (u UserServer) ChangePassword(ctx context.Context, req *proto.ChangePasswordRequest) (*proto.Empty, error) {
	var user model.User

	result := global.DB.Where(model.User{BaseModel: model.BaseModel{ID: req.Id}}).Find(&user)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效用户ID")
	}

	passwordInfo := strings.Split(user.Password, "$")
	check := password.Verify(req.GetOldPassword(), passwordInfo[2], passwordInfo[3], &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha256.New,
	})

	if !check {
		return nil, status.Errorf(codes.InvalidArgument, "原密码错误")
	}

	salt, encodedPwd := password.Encode(req.GetNewPassword(), &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha256.New,
	})

	result = global.DB.Where(model.User{BaseModel: model.BaseModel{ID: req.Id}}).Updates(&model.User{Password: fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)})
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &proto.Empty{}, nil
}
