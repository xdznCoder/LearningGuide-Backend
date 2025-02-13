package handler

import (
	"LearningGuide/file_srv/global"
	"LearningGuide/file_srv/model"
	FileProto "LearningGuide/file_srv/proto/.FileProto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f FileServer) GetNounDetail(ctx context.Context, req *FileProto.NounDetailRequest) (*FileProto.NounInfoResponse, error) {
	var noun model.Noun

	result := global.DB.Model(&model.Noun{}).Where(&model.Noun{BaseModel: model.BaseModel{ID: req.Id}}).Find(&noun)

	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "无效名词ID")
	}

	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	return &FileProto.NounInfoResponse{
		Id:       noun.ID,
		Name:     noun.Name,
		Content:  noun.Content,
		CourseId: noun.CourseId,
	}, nil
}

func (f FileServer) NewNoun(ctx context.Context, req *FileProto.NewNounRequest) (*FileProto.NewNounResponse, error) {
	var noun model.Noun

	result := global.DB.Where(&model.Noun{}).Where(&model.Noun{
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

	result = global.DB.Save(&newNoun)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "保存名词信息出错: %v", result.Error)
	}

	return &FileProto.NewNounResponse{Id: newNoun.ID}, nil
}

func (f FileServer) NounList(ctx context.Context, req *FileProto.NounListRequest) (*FileProto.NounListResponse, error) {
	if req.CourseId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "无效课程ID")
	}

	filter := global.DB.Model(&model.Noun{})

	if req.Name != "" {
		filter = filter.Where("name LIKE ?", "%"+req.Name+"%")
	}

	var count int64

	var nouns []model.Noun

	filter.Where(&model.Noun{CourseId: req.CourseId}).Count(&count)

	result := filter.Model(&model.Noun{}).Scopes(Paginate(int(req.PageNum), int(req.PageSize))).
		Where(&model.Noun{CourseId: req.CourseId}).Find(&nouns)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取名词列表出错: %v", result.Error)
	}

	var respList []*FileProto.NounInfoResponse

	for _, v := range nouns {
		respList = append(respList, &FileProto.NounInfoResponse{
			Id:       v.ID,
			Name:     v.Name,
			Content:  v.Content,
			CourseId: v.CourseId,
		})
	}

	return &FileProto.NounListResponse{
		Total: int32(count),
		Data:  respList,
	}, nil
}

func (f FileServer) DeleteNoun(ctx context.Context, req *FileProto.DeleteNounRequest) (*FileProto.Empty, error) {
	result := global.DB.Model(&model.Noun{}).Delete(&model.Noun{BaseModel: model.BaseModel{ID: req.Id}})
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效名词ID")
	}

	return &FileProto.Empty{}, nil
}
