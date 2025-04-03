package forms

type NewExerciseForms struct {
	CourseId int32 `json:"course_id" binding:"required,min=0"`
}

type UpdateExerciseForms struct {
	IsRight string `json:"is_right" binding:"required,oneof=true false"`
}
