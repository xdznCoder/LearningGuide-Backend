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

type FileListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileListLogic {
	return &FileListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FileListLogic) FileList(req *proto.FileFilterRequest) (*proto.FileListResponse, error) {
	filter := l.svcCtx.DB.Model(&model.File{})
	var files []model.File

	if req.FileName != "" {
		filter = filter.Where("file_name LIKE ?", "%"+req.FileName+"%")
	}

	if req.UserId != 0 {
		filter = filter.Where(&model.File{UserId: req.UserId})
	}

	if req.CourseId != 0 {
		filter = filter.Where(&model.File{CourseId: req.CourseId})
	}

	if req.FileType != "" {
		filter = filter.Where(&model.File{FileType: req.FileType})
	}

	var count int64
	filter.Count(&count)

	result := filter.Scopes(model.Paginate(int(req.PageNum), int(req.PageSize))).Find(&files)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "查询文件列表出错: %v", result.Error)
	}

	var respList []*proto.FileInfoResponse

	for _, v := range files {
		respList = append(respList, &proto.FileInfoResponse{
			Id:       v.ID,
			FileName: v.FileName,
			FileType: v.FileType,
			FileSize: v.FileSize,
			OssUrl:   v.OssUrl,
			Desc:     v.Desc,
			UserId:   v.UserId,
			CourseId: v.CourseId,
		})
	}

	return &proto.FileListResponse{
		Total: int32(count),
		Data:  respList,
	}, nil
}
