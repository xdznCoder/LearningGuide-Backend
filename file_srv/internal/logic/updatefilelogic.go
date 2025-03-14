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

type UpdateFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFileLogic {
	return &UpdateFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateFileLogic) UpdateFile(req *proto.UpdateFileRequest) (*proto.Empty, error) {
	file := model.File{
		Desc:    req.Desc,
		MindMap: req.MindMap,
	}

	result := l.svcCtx.DB.Model(&model.File{}).Where(model.File{BaseModel: model.BaseModel{ID: req.Id}}).Updates(&file)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效文件ID")
	}

	return &proto.Empty{}, nil
}
