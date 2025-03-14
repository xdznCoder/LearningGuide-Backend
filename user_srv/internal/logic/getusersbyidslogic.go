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

type GetUsersByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUsersByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersByIdsLogic {
	return &GetUsersByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUsersByIdsLogic) GetUsersByIds(in *userProto.IdsRequest) (*userProto.UserListResponse, error) {
	var users []model.User

	result := l.svcCtx.DB.Model(&model.User{}).Where("id IN (?)", in.Ids).Find(&users)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	var respList []*userProto.UserInfoResponse

	for _, user := range users {
		respList = append(respList, ModelToResponse(user))
	}

	return &userProto.UserListResponse{
		Total: int32(result.RowsAffected),
		Data:  respList,
	}, nil
}
