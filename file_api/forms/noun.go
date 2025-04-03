package forms

type NewNounForm struct {
	Name     string `json:"name" binding:"required"`
	CourseId int32  `json:"course_id" binding:"required,min=0"`
}
