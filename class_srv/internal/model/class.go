package model

type Course struct {
	BaseModel
	UserId      int32   `type:"type:int;not null" json:"user_id"`
	Name        string  `gorm:"type:varchar(100);not null" json:"name"`
	Type        string  `gorm:"type:varchar(100);not null" json:"type"`
	CourseSn    string  `gorm:"type:varchar(100);not null;unique" json:"course_sn"`
	Image       string  `gorm:"type:varchar(200);not null" json:"image"`
	Teacher     string  `gorm:"type:varchar(100)" json:"teacher"`
	Credit      float32 `gorm:"type:float;not null" json:"credit"`
	LessonTotal int32   `gorm:"type:int" json:"lesson_total"`
	Desc        string  `gorm:"type:text" json:"desc"`
	Term        int     `gorm:"type:tinyint;not null" json:"term"`
}

type Lesson struct {
	BaseModel
	CourseId  int32 `gorm:"type:int,not null" json:"course_id"`
	Course    Course
	UserId    int32  `type:"type:int;not null" json:"user_id"`
	Term      int    `gorm:"type:tinyint;not null" json:"term"`
	WeekNum   int    `gorm:"type:int;not null" json:"week_num"`
	DayOfWeek int    `gorm:"type:int;not null" json:"day_of_week"`
	LessonNum int    `gorm:"type:int;not null" json:"lesson_num"`
	Begin     string `gorm:"type:varchar(100);not null" json:"begin"`
	End       string `gorm:"type:varchar(100);not null" json:"end"`
}

func (Lesson) TableName() string {
	return "lesson"
}

func (Course) TableName() string {
	return "course"
}
