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

type CreateFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFileLogic {
	return &CreateFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateFileLogic) CreateFile(req *proto.CreateFileRequest) (*proto.CreateFileResponse, error) {
	var former model.File

	file := model.File{
		FileName: req.FileName,
		FileType: req.FileType,
		FileSize: req.FileSize,
		OssUrl:   req.OssUrl,
		Desc:     req.Desc,
		UserId:   req.UserId,
		CourseId: req.CourseId,
	}

	result := l.svcCtx.DB.Model(&model.File{}).
		Where(&model.File{FileName: req.GetFileName(), UserId: req.UserId, CourseId: req.CourseId}).
		Find(&former)
	if result.RowsAffected != 0 {
		file.ID = former.ID
		file.CreatedAt = former.CreatedAt
	}

	result = l.svcCtx.DB.Save(&file)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "保存文件失败: %v", result.Error)
	}

	return &proto.CreateFileResponse{Id: file.ID}, nil
}
