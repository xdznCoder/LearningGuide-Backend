package handler

import (
	"LearningGuide/file_srv/global"
	"LearningGuide/file_srv/model"
	FileProto "LearningGuide/file_srv/proto/.FileProto"
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/random"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (f FileServer) CreateSession(ctx context.Context, req *FileProto.CreateSessionRequest) (*FileProto.CreateSessionResponse, error) {
	uuid := generateSessionUuid(req.CourseId)

	session := model.Session{
		Uuid:     uuid,
		CourseId: req.CourseId,
	}

	err := global.DB.Model(model.Session{}).Create(&session).Error

	if err != nil {
		zap.S().Errorf("CreateSession err: %v", err)
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}

	return &FileProto.CreateSessionResponse{Id: session.ID}, nil
}

func (f FileServer) SessionList(ctx context.Context, req *FileProto.SessionListRequest) (*FileProto.SessionListResponse, error) {
	var sessions []model.Session

	var count int64

	global.DB.Model(model.Session{}).Where(&model.Session{CourseId: req.CourseId}).Count(&count)

	result := global.DB.Model(model.Session{}).Scopes(Paginate(int(req.PageNum), int(req.PageSize))).Where(&model.Session{
		CourseId: req.CourseId,
	}).Find(&sessions)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效课程ID")
	}

	var respList []*FileProto.SessionInfoResponse

	for _, session := range sessions {
		respList = append(respList, &FileProto.SessionInfoResponse{
			Uuid: session.Uuid,
			Id:   session.ID,
		})
	}

	return &FileProto.SessionListResponse{
		Total: int32(count),
		Data:  respList,
	}, nil
}

func (f FileServer) DeleteSession(ctx context.Context, req *FileProto.DeleteSessionRequest) (*FileProto.Empty, error) {
	var messages []model.Message

	tx := global.DB.Begin()

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

	return &FileProto.Empty{}, nil
}

func (f FileServer) NewMessage(ctx context.Context, req *FileProto.NewMessageRequest) (*FileProto.NewMessageResponse, error) {
	var session model.Session

	result := global.DB.Model(&model.Session{}).Where(&model.Session{BaseModel: model.BaseModel{ID: req.SessionId}}).Find(&session)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效会话ID")
	}

	message := model.Message{
		Content:   req.Content,
		SessionID: req.SessionId,
		Type:      int(req.Type),
		Speaker:   req.Speaker,
	}

	tx := global.DB.Begin()

	err := tx.Model(&model.Message{}).Create(&message).Error

	if err != nil {
		tx.Rollback()
		zap.S().Errorf("NewMessage err: %v", err)
		return nil, status.Errorf(codes.Internal, "创建消息失败: %v", err)
	}

	tx.Commit()

	return &FileProto.NewMessageResponse{Id: message.ID}, nil
}

func (f FileServer) MessageList(ctx context.Context, req *FileProto.MessageListRequest) (*FileProto.MessageListResponse, error) {
	var session model.Session

	result := global.DB.Model(&model.Session{}).Where(&model.Session{BaseModel: model.BaseModel{ID: req.SessionId}}).Find(&session)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效会话ID")
	}

	var messages []model.Message

	var count int64

	global.DB.Model(&model.Message{}).Scopes(Paginate(int(req.PageNum), int(req.PageSize))).
		Where(&model.Message{SessionID: req.SessionId}).
		Count(&count)

	result = global.DB.Order("add_time DESC").Model(&model.Message{}).Scopes(Paginate(int(req.PageNum), int(req.PageSize))).
		Where(&model.Message{SessionID: req.SessionId}).
		Find(&messages)
	if result.Error != nil {
		zap.S().Errorf("MessageList err: %v", result.Error)
		return nil, status.Errorf(codes.Internal, "获取消息列表失败: %v", result.Error)
	}

	var respList []*FileProto.MessageInfoResponse

	for i := len(messages) - 1; i >= 0; i-- {
		respList = append(respList, &FileProto.MessageInfoResponse{
			Id:        messages[i].ID,
			Content:   messages[i].Content,
			SessionId: messages[i].SessionID,
			Type:      int32(messages[i].Type),
			Speaker:   messages[i].Speaker,
		})
	}

	return &FileProto.MessageListResponse{
		Total: int32(count),
		Data:  respList,
	}, nil
}

func generateSessionUuid(courseId int32) string {
	t := time.Now()
	return fmt.Sprintf("%d%d%d%d%d%d%d%d",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Nanosecond(),
		courseId,
		random.RandInt(10, 99),
	)
}
