package model

type Summary struct {
	BaseModel

	WeekID       string `json:"week_id" gorm:"column:week_id;type:varchar(50);not null"`
	CourseID     int32  `json:"course_id" gorm:"type:int;not null"`
	ExerciseDone string `json:"exercise_done" gorm:"column:exercise_done;type:varchar(50);not null"`
	AccuracyRate string `json:"accuracy_rate" gorm:"column:accuracy_rate;type:varchar(50);not null"`
	SessionNum   int32  `json:"session_num" gorm:"column:session_num;type:int;not null"`
	MessageNum   int32  `json:"message_num" gorm:"column:message_num;type:int;not null"`
	NounNum      int32  `json:"noun_num" gorm:"column:noun_num;type:int;not null"`
}

func (Summary) TableName() string {
	return "summary"
}
