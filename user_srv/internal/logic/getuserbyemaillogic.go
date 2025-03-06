package logic

import (
	"LearningGuide/user_srv/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"LearningGuide/user_srv/internal/svc"
	"LearningGuide/user_srv/userProto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByEmailLogic {
	return &GetUserByEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByEmailLogic) GetUserByEmail(in *userProto.EmailRequest) (*userProto.UserInfoResponse, error) {
	var user model.User
	result := l.svcCtx.DB.Where(&model.User{Email: in.Email}).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效用户邮箱")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return ModelToResponse(user), nil
}
