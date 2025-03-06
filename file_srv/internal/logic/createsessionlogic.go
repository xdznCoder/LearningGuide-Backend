package logic

import (
	"LearningGuide/file_srv/internal/model"
	"context"
	"fmt"
	"github.com/duke-git/lancet/v2/random"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	proto "LearningGuide/file_srv/.FileProto"
	"LearningGuide/file_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSessionLogic) CreateSession(req *proto.CreateSessionRequest) (*proto.CreateSessionResponse, error) {
	uuid := generateSessionUuid(req.CourseId)

	session := model.Session{
		Uuid:     uuid,
		CourseId: req.CourseId,
	}

	err := l.svcCtx.DB.Model(model.Session{}).Create(&session).Error

	if err != nil {
		zap.S().Errorf("CreateSession err: %v", err)
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}

	return &proto.CreateSessionResponse{Id: session.ID}, nil
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
