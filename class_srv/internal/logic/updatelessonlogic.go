package logic

import (
	"LearningGuide/class_srv/internal/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "LearningGuide/class_srv/.ClassProto"
	"LearningGuide/class_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLessonLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLessonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLessonLogic {
	return &UpdateLessonLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLessonLogic) UpdateLesson(in *proto.UpdateLessonRequest) (*proto.Empty, error) {
	if !beginEarlierThanEnd(in.Begin, in.End) {
		return nil, status.Errorf(codes.InvalidArgument, "无效的课时始末时间")
	}

	lesson := model.Lesson{
		Begin: in.Begin,
		End:   in.End,
	}

	result := l.svcCtx.DB.Model(&model.Lesson{}).Where(&model.Lesson{BaseModel: model.BaseModel{ID: in.Id}}).Updates(&lesson)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效课时ID")
	}

	return &proto.Empty{}, nil
}
