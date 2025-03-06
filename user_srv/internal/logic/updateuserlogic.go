package logic

import (
	"LearningGuide/user_srv/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	"LearningGuide/user_srv/internal/svc"
	"LearningGuide/user_srv/userProto"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *userProto.UpdateUserInfo) (*userProto.Empty, error) {
	var user model.User
	result := l.svcCtx.DB.Where(&model.User{BaseModel: model.BaseModel{ID: in.Id}}).Find(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效用户ID")
	}
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	birthDay := time.Unix(int64(in.GetBirthDay()), 0)
	user.NickName = in.GetNickName()
	user.Birthday = &birthDay
	user.Gender = in.GetGender()
	user.Desc = in.GetDesc()
	user.Image = in.GetImage()

	result = l.svcCtx.DB.Updates(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &userProto.Empty{}, nil
}
