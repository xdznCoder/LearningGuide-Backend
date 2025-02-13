package forms

type GenerateExerciseForms struct {
	CourseId int32   `json:"course_id" binding:"required,min=0"`
	FileIds  []int32 `json:"file_ids"`
}

type NewExerciseForms struct {
	CourseId int32  `json:"course_id" binding:"required,min=0"`
	ResultId string `json:"result_id" binding:"required"`
}

type UpdateExerciseForms struct {
	IsRight string `json:"is_right" binding:"required,oneof=true false"`
}
