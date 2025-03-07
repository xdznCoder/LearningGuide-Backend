package logic

import (
	"LearningGuide/post_srv/internal/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "LearningGuide/post_srv/.PostProto"
	"LearningGuide/post_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateNoticeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateNoticeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNoticeListLogic {
	return &UpdateNoticeListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateNoticeListLogic) UpdateNoticeList(in *proto.NoticeUpdateRequest) (*proto.Empty, error) {
	result := l.svcCtx.DB.Model(&model.Notice{}).
		Where(&model.Notice{BaseModel: model.BaseModel{ID: in.Id}}).
		Update("is_read", true)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的通知ID")
	}

	return &proto.Empty{}, nil
}
