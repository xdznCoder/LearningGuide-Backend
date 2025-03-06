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

type DeleteNounLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNounLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNounLogic {
	return &DeleteNounLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteNounLogic) DeleteNoun(req *proto.DeleteNounRequest) (*proto.Empty, error) {
	result := l.svcCtx.DB.Model(&model.Noun{}).Delete(&model.Noun{BaseModel: model.BaseModel{ID: req.Id}})
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效名词ID")
	}

	return &proto.Empty{}, nil
}
