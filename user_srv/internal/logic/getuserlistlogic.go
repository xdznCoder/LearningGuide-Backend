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

type GetUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserListLogic) GetUserList(in *userProto.PageInfo) (*userProto.UserListResponse, error) {
	var users []model.User
	userListResponse := make([]*userProto.UserInfoResponse, 0)

	result := l.svcCtx.DB.Find(&users)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	err := l.svcCtx.DB.Scopes(model.Paginate(int(in.PageNum), int(in.PageSize))).Find(&users).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	for _, user := range users {
		userListResponse = append(userListResponse, ModelToResponse(user))
	}

	return &userProto.UserListResponse{
		Total: int32(result.RowsAffected),
		Data:  userListResponse,
	}, nil
}
