package logic

import (
	"LearningGuide/user_srv/model"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"

	"LearningGuide/user_srv/internal/svc"
	"LearningGuide/user_srv/userProto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangePasswordLogic) ChangePassword(in *userProto.ChangePasswordRequest) (*userProto.Empty, error) {
	var user model.User

	result := l.svcCtx.DB.Where(model.User{BaseModel: model.BaseModel{ID: in.Id}}).Find(&user)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效用户ID")
	}

	passwordInfo := strings.Split(user.Password, "$")
	check := password.Verify(in.GetOldPassword(), passwordInfo[2], passwordInfo[3], &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha256.New,
	})

	if !check {
		return nil, status.Errorf(codes.InvalidArgument, "原密码错误")
	}

	salt, encodedPwd := password.Encode(in.GetNewPassword(), &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha256.New,
	})

	result = l.svcCtx.DB.Where(model.User{BaseModel: model.BaseModel{ID: in.Id}}).Updates(
		&model.User{
			Password: fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd),
		})
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &userProto.Empty{}, nil
}
