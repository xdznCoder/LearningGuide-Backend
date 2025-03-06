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

type DeleteFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFileLogic) DeleteFile(req *proto.DeleteFileRequest) (*proto.Empty, error) {
	result := l.svcCtx.DB.Model(&model.File{}).Delete(&model.File{BaseModel: model.BaseModel{ID: req.Id}})
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效文件ID")
	}

	return &proto.Empty{}, nil

}
