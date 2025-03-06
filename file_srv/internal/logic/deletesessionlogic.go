package logic

import (
	"LearningGuide/file_srv/internal/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "LearningGuide/file_srv/.FileProto"
	"LearningGuide/file_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSessionLogic {
	return &DeleteSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSessionLogic) DeleteSession(req *proto.DeleteSessionRequest) (*proto.Empty, error) {
	var messages []model.Message

	tx := l.svcCtx.DB.Begin()

	result := tx.Delete(&model.Session{BaseModel: model.BaseModel{ID: req.Id}})
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "无效会话ID")
	}

	err := tx.Model(&model.Message{}).Unscoped().Where(&model.Message{SessionID: req.Id}).Delete(&messages).Error

	if err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "删除消息记录失败: %v", err)
	}

	tx.Commit()

	return &proto.Empty{}, nil
}
