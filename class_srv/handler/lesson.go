package handler

import (
	"LearningGuide/class_srv/global"
	"LearningGuide/class_srv/model"
	proto "LearningGuide/class_srv/proto/.ClassProto"
	"context"
	"github.com/duke-git/lancet/v2/slice"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"strconv"
	"strings"
)

func (c ClassServer) GetLessonList(ctx context.Context, req *proto.LessonFilterRequest) (*proto.LessonListResponse, error) {
	var RespList []*proto.LessonInfoResponse

	filter := global.DB.Preload("Course").Model(&model.Lesson{})

	if req.UserId != 0 {
		filter = filter.Where(&model.Lesson{UserId: req.UserId})
	}

	if req.DayOfWeek != 0 {
		filter = filter.Where(&model.Lesson{DayOfWeek: int(req.DayOfWeek)})
	}

	if req.WeekNum != 0 {
		filter = filter.Where(&model.Lesson{WeekNum: int(req.WeekNum)})
	}

	if req.CourseId != 0 {
		filter = filter.Where(&model.Lesson{CourseId: req.CourseId})
	}

	var lessons []model.Lesson

	result := filter.Find(&lessons)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	for _, v := range lessons {
		RespList = append(RespList, &proto.LessonInfoResponse{
			Id:        v.ID,
			UserId:    v.UserId,
			CourseId:  v.CourseId,
			WeekNum:   int32(v.WeekNum),
			DayOfWeek: int32(v.DayOfWeek),
			LessonNum: int32(v.LessonNum),
			Begin:     v.Begin,
			End:       v.End,
			Course: &proto.CourseInfoResponse{
				Id:       v.CourseId,
				Name:     v.Course.Name,
				Type:     v.Course.Type,
				CourseSn: v.Course.CourseSn,
				Image:    v.Course.Image,
				Teacher:  v.Course.Teacher,
				Credit:   v.Course.Credit,
			},
		})
	}

	return &proto.LessonListResponse{
		Total: int32(result.RowsAffected),
		Data:  RespList,
	}, nil
}

func (c ClassServer) CreateLesson(ctx context.Context, req *proto.CreateLessonRequest) (*proto.CreateLessonResponse, error) {
	if !beginEarlierThanEnd(req.Begin, req.End) {
		return nil, status.Errorf(codes.InvalidArgument, "无效的课时始末时间")
	}

	tx := global.DB.Begin()

	var course model.Course
	result := tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: req.CourseId}}).Find(&course)

	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "无效课程ID")
	}

	if course.UserId != req.UserId {
		tx.Rollback()
		return nil, status.Errorf(codes.InvalidArgument, "无效用户ID")
	}

	var former model.Lesson
	result = tx.Model(&model.Lesson{}).Where(&model.Lesson{
		UserId:    req.UserId,
		Term:      course.Term,
		WeekNum:   int(req.WeekNum),
		DayOfWeek: int(req.DayOfWeek),
		LessonNum: int(req.LessonNum),
	}).Find(&former)
	if result.RowsAffected != 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.AlreadyExists, "该段时间已有其他课程")
	}

	lesson := model.Lesson{
		CourseId:  req.CourseId,
		UserId:    req.UserId,
		Term:      course.Term,
		WeekNum:   int(req.WeekNum),
		DayOfWeek: int(req.DayOfWeek),
		LessonNum: int(req.LessonNum),
		Begin:     req.Begin,
		End:       req.End,
	}

	result = tx.Create(&lesson)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	result = tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: req.CourseId}}).Update("lesson_total", course.LessonTotal+1)

	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	tx.Commit()

	return &proto.CreateLessonResponse{Id: lesson.ID}, nil
}

func (c ClassServer) CreateLessonInBatch(ctx context.Context, req *proto.CreateLessonBatchRequest) (*proto.CreateLessonBatchResponse, error) {
	var ids []int32

	var lessons []model.Lesson

	if !beginEarlierThanEnd(req.Begin, req.End) {
		return nil, status.Errorf(codes.InvalidArgument, "无效的课时始末时间")
	}

	if req.EndWeek < req.BeginWeek {
		return nil, status.Errorf(codes.InvalidArgument, "无效的课程始末时间")
	}

	tx := global.DB.Begin()

	var course model.Course
	result := tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: req.CourseId}}).Find(&course)

	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "无效课程ID")
	}

	if course.UserId != req.UserId {
		tx.Rollback()
		return nil, status.Errorf(codes.InvalidArgument, "无效用户ID")
	}

	var former []model.Lesson

	tx.Model(&model.Lesson{}).Where(model.Lesson{
		UserId:    req.UserId,
		CourseId:  req.CourseId,
		Term:      course.Term,
		DayOfWeek: int(req.DayOfWeek),
		LessonNum: int(req.LessonNum),
	}).Find(&former)

	for _, v := range former {
		if int32(v.WeekNum) <= req.EndWeek && int32(v.WeekNum) >= req.BeginWeek {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "存在课时时间冲突")
		}
	}

	for i := req.BeginWeek; i <= req.EndWeek; i++ {
		lessons = append(lessons, model.Lesson{
			CourseId:  req.CourseId,
			UserId:    req.UserId,
			Term:      course.Term,
			WeekNum:   int(i),
			DayOfWeek: int(req.DayOfWeek),
			LessonNum: int(req.LessonNum),
			Begin:     req.Begin,
			End:       req.End,
		})
	}

	result = tx.Model(model.Lesson{}).Create(lessons)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.InvalidArgument, result.Error.Error())
	}

	result = tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: req.CourseId}}).
		Update("lesson_total", course.LessonTotal+req.EndWeek-req.BeginWeek+1)

	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	tx.Commit()

	for _, v := range lessons {
		ids = append(ids, v.ID)
	}

	return &proto.CreateLessonBatchResponse{Ids: ids}, nil
}

func (c ClassServer) UpdateLesson(ctx context.Context, req *proto.UpdateLessonRequest) (*proto.Empty, error) {
	if !beginEarlierThanEnd(req.Begin, req.End) {
		return nil, status.Errorf(codes.InvalidArgument, "无效的课时始末时间")
	}

	lesson := model.Lesson{
		Begin: req.Begin,
		End:   req.End,
	}

	result := global.DB.Model(&model.Lesson{}).Where(&model.Lesson{BaseModel: model.BaseModel{ID: req.Id}}).Updates(&lesson)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效课时ID")
	}

	return &proto.Empty{}, nil
}

func (c ClassServer) DeleteLesson(ctx context.Context, req *proto.DeleteLessonRequest) (*proto.Empty, error) {
	var lesson model.Lesson

	tx := global.DB.Begin()

	result := tx.Model(&model.Lesson{}).Where(&model.Lesson{BaseModel: model.BaseModel{ID: req.Id}}).Find(&lesson)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "无效课时ID")
	}

	var course model.Course
	result = tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: lesson.CourseId}}).Find(&course)

	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	result = tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: lesson.CourseId}}).Update("lesson_total", course.LessonTotal-1)

	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	result = tx.Model(&model.Lesson{}).Delete(&model.Lesson{BaseModel: model.BaseModel{ID: req.Id}})
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "无效课时ID")
	}

	tx.Commit()

	return &proto.Empty{}, nil
}

func (c ClassServer) GetLessonDetail(ctx context.Context, req *proto.LessonDetailRequest) (*proto.LessonInfoResponse, error) {

	var lesson model.Lesson

	result := global.DB.Preload("Course").Where(&model.Lesson{BaseModel: model.BaseModel{ID: req.Id}}).Find(&lesson)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "无效课时ID")
	}

	return &proto.LessonInfoResponse{
		Id:        lesson.ID,
		CourseId:  lesson.CourseId,
		WeekNum:   int32(lesson.WeekNum),
		DayOfWeek: int32(lesson.DayOfWeek),
		LessonNum: int32(lesson.LessonNum),
		Begin:     lesson.Begin,
		End:       lesson.End,
		Course: &proto.CourseInfoResponse{
			Id:          lesson.Course.ID,
			Name:        lesson.Course.Name,
			Type:        lesson.Course.Type,
			CourseSn:    lesson.Course.CourseSn,
			Term:        int32(lesson.Course.Term),
			LessonTotal: lesson.Course.LessonTotal,
			Desc:        lesson.Course.Desc,
			Image:       lesson.Course.Image,
			Teacher:     lesson.Course.Teacher,
			Credit:      lesson.Course.Credit,
			UserId:      lesson.Course.UserId,
		},
		Term:   int32(lesson.Term),
		UserId: lesson.UserId,
	}, nil
}

func (c ClassServer) DeleteLessonInBatch(ctx context.Context, req *proto.DeleteLessonInBatchRequest) (*proto.Empty, error) {
	req.Ids = slice.Unique(req.Ids)

	tx := global.DB.Begin()

	var course model.Course
	result := tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: req.CourseId}}).Find(&course)

	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "无效课程ID")
	}

	result = tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: req.CourseId}}).Update("lesson_total", int(course.LessonTotal)-len(req.Ids))

	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	var conditions [][]interface{}
	for _, id := range req.Ids {
		conditions = append(conditions, []interface{}{id, req.UserId, req.CourseId})
	}

	result = tx.Model(&model.Lesson{}).Where("(id, user_id, course_id) IN (?)", conditions).Delete(&model.Lesson{})
	if result.RowsAffected < int64(len(req.Ids)) {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "存在无效课时ID")
	}

	tx.Commit()
	return &proto.Empty{}, nil
}

func beginEarlierThanEnd(begin string, end string) bool {
	timeRegex := regexp.MustCompile(`^([01]?[0-9]|2[0-3]):([0-5][0-9])$`)

	if !timeRegex.MatchString(begin) || !timeRegex.MatchString(end) {
		return false
	}

	bs := strings.Split(begin, ":")
	es := strings.Split(end, ":")
	bh, err := strconv.Atoi(bs[0])
	eh, err := strconv.Atoi(es[0])
	bm, err := strconv.Atoi(bs[1])
	em, err := strconv.Atoi(es[1])
	if err != nil {
		return false
	}
	return bh*60+bm < eh*60+em
}
