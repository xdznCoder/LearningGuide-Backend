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

type MessageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageListLogic {
	return &MessageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MessageListLogic) MessageList(req *proto.MessageListRequest) (*proto.MessageListResponse, error) {
	var session model.Session

	result := l.svcCtx.DB.Model(&model.Session{}).Where(&model.Session{BaseModel: model.BaseModel{ID: req.SessionId}}).Find(&session)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效会话ID")
	}

	var messages []model.Message

	var count int64

	l.svcCtx.DB.Model(&model.Message{}).Scopes(model.Paginate(int(req.PageNum), int(req.PageSize))).
		Where(&model.Message{SessionID: req.SessionId}).
		Count(&count)

	result = l.svcCtx.DB.Order("add_time DESC").Model(&model.Message{}).Scopes(model.Paginate(int(req.PageNum), int(req.PageSize))).
		Where(&model.Message{SessionID: req.SessionId}).
		Find(&messages)
	if result.Error != nil {
		zap.S().Errorf("MessageList err: %v", result.Error)
		return nil, status.Errorf(codes.Internal, "获取消息列表失败: %v", result.Error)
	}

	var respList []*proto.MessageInfoResponse

	for i := len(messages) - 1; i >= 0; i-- {
		respList = append(respList, &proto.MessageInfoResponse{
			Id:        messages[i].ID,
			Content:   messages[i].Content,
			SessionId: messages[i].SessionID,
			Type:      int32(messages[i].Type),
			Speaker:   messages[i].Speaker,
		})
	}

	return &proto.MessageListResponse{
		Total: int32(count),
		Data:  respList,
	}, nil
}
