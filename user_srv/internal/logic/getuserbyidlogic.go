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

type GetUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByIdLogic) GetUserById(in *userProto.IdRequest) (*userProto.UserInfoResponse, error) {
	var user model.User
	result := l.svcCtx.DB.Where(&model.User{BaseModel: model.BaseModel{ID: in.Id}}).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效用户ID")
	}
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return ModelToResponse(user), nil
}
