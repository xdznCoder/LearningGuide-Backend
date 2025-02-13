package handler

import (
	"LearningGuide/file_srv/global"
	"LearningGuide/file_srv/model"
	proto "LearningGuide/file_srv/proto/.FileProto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type FileServer struct {
	proto.UnimplementedFileServer
}

func (f FileServer) CreateFile(ctx context.Context, req *proto.CreateFileRequest) (*proto.CreateFileResponse, error) {
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

	result := global.DB.Model(&model.File{}).
		Where(&model.File{FileName: req.GetFileName(), UserId: req.UserId, CourseId: req.CourseId}).
		Find(&former)
	if result.RowsAffected != 0 {
		file.ID = former.ID
		file.CreatedAt = former.CreatedAt
	}

	result = global.DB.Save(&file)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "保存文件失败: %v", result.Error)
	}

	return &proto.CreateFileResponse{Id: file.ID}, nil
}

func (f FileServer) GetFileDetail(ctx context.Context, req *proto.FileDetailRequest) (*proto.FileInfoResponse, error) {
	var file model.File

	result := global.DB.Model(&model.File{}).
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
	}, nil
}

func (f FileServer) FileList(ctx context.Context, req *proto.FileFilterRequest) (*proto.FileListResponse, error) {
	filter := global.DB.Model(&model.File{})
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

	result := filter.Scopes(Paginate(int(req.PageNum), int(req.PageSize))).Find(&files)
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

func (f FileServer) UpdateFile(ctx context.Context, req *proto.UpdateFileRequest) (*proto.Empty, error) {
	file := model.File{
		Desc: req.Desc,
	}

	result := global.DB.Model(&model.File{}).Where(model.File{BaseModel: model.BaseModel{ID: req.Id}}).Updates(&file)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效文件ID")
	}

	return &proto.Empty{}, nil
}

func (f FileServer) DeleteFile(ctx context.Context, req *proto.DeleteFileRequest) (*proto.Empty, error) {
	result := global.DB.Model(&model.File{}).Delete(&model.File{BaseModel: model.BaseModel{ID: req.Id}})
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效文件ID")
	}

	return &proto.Empty{}, nil
}

func Paginate(pageNum, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNum == 0 {
			pageNum = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
