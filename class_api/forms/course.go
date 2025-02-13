package forms

type CreateCourseForm struct {
	Name    string  `json:"name" binding:"required,min=1,max=20"`
	Type    string  `json:"type" binding:"required"`
	Term    int     `json:"term" binding:"required,oneof=1 2 3 4 5 6 7 8"`
	Desc    string  `json:"desc" binding:"required"`
	Image   string  `json:"image" binding:"required,url"`
	Credit  float32 `json:"credit" binding:"required,min=0"`
	Teacher string  `json:"teacher" binding:"required"`
	UserId  int     `json:"userId" binding:"required"`
}

type UpdateCourseForm struct {
	Name    string  `json:"name"`
	Desc    string  `json:"desc"`
	Image   string  `json:"image"`
	Teacher string  `json:"teacher"`
	Credit  float32 `json:"credit" binding:"min=0"`
}

type CreateLessonForm struct {
	CourseId  int    `json:"course_id" binding:"required"`
	WeekNum   int    `json:"week_num" binding:"required"`
	DayOfWeek int    `json:"day_of_week" binding:"required"`
	LessonNum int    `json:"lesson_num" binding:"required"`
	Begin     string `json:"begin" binding:"required,min=5,max=5"`
	End       string `json:"end" binding:"required,min=5,max=5"`
	UserId    int    `json:"user_id" binding:"required"`
}

type UpdateLessonForm struct {
	Begin string `json:"begin"`
	End   string `json:"end"`
}

type CreateLessonFormInBatch struct {
	CourseId  int    `json:"course_id" binding:"required"`
	UserId    int    `json:"user_id" binding:"required"`
	BeginWeek int    `json:"begin_week" binding:"required,min=0"`
	EndWeek   int    `json:"end_week" binding:"required,min=0"`
	DayOfWeek int    `json:"day_of_Week" binding:"required,min=0"`
	LessonNum int    `json:"lesson_num" binding:"required,min=0"`
	Begin     string `json:"begin" binding:"required"`
	End       string `json:"end" binding:"required"`
}

type DeleteLessonFormInBatch struct {
	CourseId int     `json:"course_id" binding:"required"`
	UserId   int     `json:"user_id" binding:"required"`
	Ids      []int32 `json:"ids" binding:"required"`
}
