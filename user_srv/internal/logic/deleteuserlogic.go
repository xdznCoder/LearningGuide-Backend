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

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteUserLogic) DeleteUser(in *userProto.DeleteUserRequest) (*userProto.Empty, error) {
	result := l.svcCtx.DB.Where(&model.User{BaseModel: model.BaseModel{ID: in.Id}}).Delete(&model.User{})
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效用户ID")
	}

	return &userProto.Empty{}, nil
}
