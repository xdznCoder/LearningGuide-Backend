package logic

import (
	"LearningGuide/class_srv/internal/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"strconv"
	"strings"

	proto "LearningGuide/class_srv/.ClassProto"
	"LearningGuide/class_srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLessonLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLessonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLessonLogic {
	return &CreateLessonLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLessonLogic) CreateLesson(in *proto.CreateLessonRequest) (*proto.CreateLessonResponse, error) {
	if !beginEarlierThanEnd(in.Begin, in.End) {
		return nil, status.Errorf(codes.InvalidArgument, "无效的课时始末时间")
	}

	tx := l.svcCtx.DB.Begin()

	var course model.Course
	result := tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: in.CourseId}}).Find(&course)

	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "无效课程ID")
	}

	if course.UserId != in.UserId {
		tx.Rollback()
		return nil, status.Errorf(codes.InvalidArgument, "无效用户ID")
	}

	var former model.Lesson
	result = tx.Model(&model.Lesson{}).Where(&model.Lesson{
		UserId:    in.UserId,
		Term:      course.Term,
		WeekNum:   int(in.WeekNum),
		DayOfWeek: int(in.DayOfWeek),
		LessonNum: int(in.LessonNum),
	}).Find(&former)
	if result.RowsAffected != 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.AlreadyExists, "该段时间已有其他课程")
	}

	lesson := model.Lesson{
		CourseId:  in.CourseId,
		UserId:    in.UserId,
		Term:      course.Term,
		WeekNum:   int(in.WeekNum),
		DayOfWeek: int(in.DayOfWeek),
		LessonNum: int(in.LessonNum),
		Begin:     in.Begin,
		End:       in.End,
	}

	result = tx.Create(&lesson)
	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	result = tx.Model(&model.Course{}).Where(&model.Course{BaseModel: model.BaseModel{ID: in.CourseId}}).Update("lesson_total", course.LessonTotal+1)

	if result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	tx.Commit()

	return &proto.CreateLessonResponse{Id: lesson.ID}, nil
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
