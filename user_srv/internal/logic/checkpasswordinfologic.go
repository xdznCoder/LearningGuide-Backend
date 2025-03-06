package logic

import (
	"context"
	"crypto/sha256"
	"github.com/anaskhan96/go-password-encoder"
	"strings"

	"LearningGuide/user_srv/internal/svc"
	"LearningGuide/user_srv/userProto"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckPasswordInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckPasswordInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPasswordInfoLogic {
	return &CheckPasswordInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckPasswordInfoLogic) CheckPasswordInfo(in *userProto.PasswordCheck) (*userProto.CheckResponse, error) {
	passwordInfo := strings.Split(in.GetEncryptedPassword(), "$")
	check := password.Verify(in.GetPassword(), passwordInfo[2], passwordInfo[3], &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha256.New,
	})

	return &userProto.CheckResponse{Success: check}, nil
}
