package logic

import (
	proto "LearningGuide/post_srv/.PostProto"
	"LearningGuide/post_srv/internal/model"
	"LearningGuide/post_srv/internal/svc"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckNoticeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckNoticeLogic {
	return &CheckNoticeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckNoticeLogic) CheckNotice(in *proto.CheckNoticeRequest) (*proto.CheckNoticeResponse, error) {
	result := l.svcCtx.DB.Model(&model.Notice{}).Where(&model.Notice{OwnerId: in.UserId}).
		Where("owner_id != user_id").
		Where("is_read = ?", 0).
		Limit(1).Find(&model.Notice{})

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return &proto.CheckNoticeResponse{NewNotices: false}, nil
	}

	return &proto.CheckNoticeResponse{NewNotices: true}, nil
}
