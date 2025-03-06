package logic

import (
	"LearningGuide/user_srv/model"
	"context"
	"crypto/sha256"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"LearningGuide/user_srv/internal/svc"
	"LearningGuide/user_srv/userProto"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUserLogic) CreateUser(in *userProto.CreateUserInfo) (*userProto.UserInfoResponse, error) {
	var formerUser model.User
	result := l.svcCtx.DB.Where(&model.User{Email: in.Email}).Find(&formerUser)
	if result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "电子邮箱已被使用")
	}

	salt, encodedPwd := password.Encode(in.GetPassword(), &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha256.New,
	})

	user := model.User{
		Email:    in.Email,
		Password: fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd),
		NickName: in.GetNickName(),
	}

	result = l.svcCtx.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建用户失败: %v", result.Error)
	}

	return ModelToResponse(user), nil
}
