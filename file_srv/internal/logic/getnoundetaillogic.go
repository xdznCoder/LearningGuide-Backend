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

type GetNounDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNounDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNounDetailLogic {
	return &GetNounDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNounDetailLogic) GetNounDetail(req *proto.NounDetailRequest) (*proto.NounInfoResponse, error) {
	var noun model.Noun

	result := l.svcCtx.DB.Model(&model.Noun{}).Where(&model.Noun{BaseModel: model.BaseModel{ID: req.Id}}).Find(&noun)

	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "无效名词ID")
	}

	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	return &proto.NounInfoResponse{
		Id:       noun.ID,
		Name:     noun.Name,
		Content:  noun.Content,
		CourseId: noun.CourseId,
	}, nil
}
