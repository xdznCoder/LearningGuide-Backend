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

type GetFileDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileDetailLogic {
	return &GetFileDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFileDetailLogic) GetFileDetail(req *proto.FileDetailRequest) (*proto.FileInfoResponse, error) {
	var file model.File

	result := l.svcCtx.DB.Model(&model.File{}).
		Where(&model.File{BaseModel: model.BaseModel{ID: req.Id}}).
		First(&file)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效文件ID")
	}

	return &proto.FileInfoResponse{
		Id:       file.ID,
		FileName: file.FileName,
		FileType: file.FileType,
		FileSize: file.FileSize,
		OssUrl:   file.OssUrl,
		Desc:     file.Desc,
		UserId:   file.UserId,
		CourseId: file.CourseId,
		MindMap:  file.MindMap,
	}, nil
}
