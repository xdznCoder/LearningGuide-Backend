package logic

import (
	"LearningGuide/file_srv/internal/model"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "LearningGuide/file_srv/.FileProto"
	"LearningGuide/file_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type NewMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNewMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewMessageLogic {
	return &NewMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NewMessageLogic) NewMessage(req *proto.NewMessageRequest) (*proto.NewMessageResponse, error) {
	var session model.Session

	result := l.svcCtx.DB.Model(&model.Session{}).Where(&model.Session{BaseModel: model.BaseModel{ID: req.SessionId}}).Find(&session)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效会话ID")
	}

	message := model.Message{
		Content:   req.Content,
		SessionID: req.SessionId,
		Type:      int(req.Type),
		Speaker:   req.Speaker,
	}

	tx := l.svcCtx.DB.Begin()

	err := tx.Model(&model.Message{}).Create(&message).Error

	if err != nil {
		tx.Rollback()
		zap.S().Errorf("NewMessage err: %v", err)
		return nil, status.Errorf(codes.Internal, "创建消息失败: %v", err)
	}

	tx.Commit()

	return &proto.NewMessageResponse{Id: message.ID}, nil
}
