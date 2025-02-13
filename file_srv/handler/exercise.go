package handler

import (
	"LearningGuide/file_srv/global"
	"LearningGuide/file_srv/model"
	FileProto "LearningGuide/file_srv/proto/.FileProto"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f FileServer) NewExercise(ctx context.Context, req *FileProto.NewExerciseRequest) (*FileProto.NewExerciseResponse, error) {
	exercise := model.Exercise{
		CourseId: req.CourseId,
		Question: req.Question,
		SectionA: req.SectionA,
		SectionB: req.SectionB,
		SectionC: req.SectionC,
		SectionD: req.SectionD,
		Answer:   req.Answer,
		Reason:   req.Reason,
	}

	result := global.DB.Model(&model.Exercise{}).Create(&exercise)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &FileProto.NewExerciseResponse{Id: exercise.ID}, nil
}

func (f FileServer) UpdateExercise(ctx context.Context, req *FileProto.UpdateExerciseRequest) (*FileProto.Empty, error) {
	fmt.Println(req.IsRight)

	switch req.IsRight {
	case "true":
	case "false":
	default:
		return nil, status.Errorf(codes.InvalidArgument, "无效的is_right参数")
	}

	result := global.DB.Model(&model.Exercise{}).
		Where(&model.Exercise{BaseModel: model.BaseModel{ID: req.Id}}).
		Update("is_right", req.IsRight)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的习题ID")
	}

	return &FileProto.Empty{}, nil
}

func (f FileServer) ExerciseList(ctx context.Context, req *FileProto.ExerciseListRequest) (*FileProto.ExerciseListResponse, error) {
	filter := global.DB.Model(&model.Exercise{})

	if req.CourseId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "无效的课程ID")
	}

	if req.Question != "" {
		filter = filter.Where("question LIKE ?", "%"+req.Question+"%")
	}

	if req.IsRight != "" {
		switch req.IsRight {
		case "true":
			filter = filter.Where(&model.Exercise{IsRight: "true"})
		case "false":
			filter = filter.Where(&model.Exercise{IsRight: "false"})
		default:
			return nil, status.Errorf(codes.InvalidArgument, "无效的is_right参数")
		}
	}

	var count int64

	var exercises []model.Exercise

	filter.Count(&count)

	result := filter.Scopes(Paginate(int(req.PageNum), int(req.PageSize))).Find(&exercises)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	var respList []*FileProto.ExerciseInfoResponse

	for _, v := range exercises {
		respList = append(respList, &FileProto.ExerciseInfoResponse{
			CourseId: v.CourseId,
			Question: v.Question,
			SectionA: v.SectionA,
			SectionB: v.SectionB,
			SectionC: v.SectionC,
			SectionD: v.SectionD,
			Answer:   v.Answer,
			Reason:   v.Reason,
			IsRight:  v.IsRight,
			Id:       v.ID,
		})
	}

	return &FileProto.ExerciseListResponse{
		Total: int32(count),
		Data:  respList,
	}, nil
}

func (f FileServer) GetExerciseDetail(ctx context.Context, req *FileProto.ExerciseDetailRequest) (*FileProto.ExerciseInfoResponse, error) {
	var exer model.Exercise

	result := global.DB.Model(&model.Exercise{}).Where(&model.Exercise{BaseModel: model.BaseModel{ID: req.Id}}).Find(&exer)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的习题ID")
	}

	return &FileProto.ExerciseInfoResponse{
		CourseId: exer.CourseId,
		Question: exer.Question,
		SectionA: exer.SectionA,
		SectionB: exer.SectionB,
		SectionC: exer.SectionC,
		SectionD: exer.SectionD,
		Answer:   exer.Answer,
		Reason:   exer.Reason,
		IsRight:  exer.IsRight,
		Id:       exer.ID,
	}, nil
}

func (f FileServer) DeleteExercise(ctx context.Context, req *FileProto.DeleteExerciseRequest) (*FileProto.Empty, error) {
	result := global.DB.Model(&model.Exercise{}).Delete(&model.Exercise{BaseModel: model.BaseModel{ID: req.Id}})

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效的课时ID")
	}

	return &FileProto.Empty{}, nil
}
