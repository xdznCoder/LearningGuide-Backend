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

type NewNounLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNewNounLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewNounLogic {
	return &NewNounLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NewNounLogic) NewNoun(req *proto.NewNounRequest) (*proto.NewNounResponse, error) {
	var noun model.Noun

	result := l.svcCtx.DB.Where(&model.Noun{}).Where(&model.Noun{
		Name:     req.Name,
		CourseId: req.CourseId,
	}).Find(&noun)

	newNoun := model.Noun{
		Name:     req.Name,
		Content:  req.Content,
		CourseId: req.CourseId,
	}

	if result.RowsAffected != 0 {
		newNoun.ID = noun.ID
	}

	result = l.svcCtx.DB.Save(&newNoun)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "保存名词信息出错: %v", result.Error)
	}

	return &proto.NewNounResponse{Id: newNoun.ID}, nil
}
