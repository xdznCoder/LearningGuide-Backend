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

type SessionListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSessionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SessionListLogic {
	return &SessionListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SessionListLogic) SessionList(req *proto.SessionListRequest) (*proto.SessionListResponse, error) {
	var sessions []model.Session

	var count int64

	l.svcCtx.DB.Model(model.Session{}).Where(&model.Session{CourseId: req.CourseId}).Count(&count)

	result := l.svcCtx.DB.Model(model.Session{}).Scopes(model.Paginate(int(req.PageNum), int(req.PageSize))).Where(&model.Session{
		CourseId: req.CourseId,
	}).Find(&sessions)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效课程ID")
	}

	var respList []*proto.SessionInfoResponse

	for _, session := range sessions {
		respList = append(respList, &proto.SessionInfoResponse{
			Uuid: session.Uuid,
			Id:   session.ID,
		})
	}

	return &proto.SessionListResponse{
		Total: int32(count),
		Data:  respList,
	}, nil
}
