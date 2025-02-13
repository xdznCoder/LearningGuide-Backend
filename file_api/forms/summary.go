package forms

type NewSummaryForms struct {
	CourseId int32  `json:"course_id" binding:"required,min=0"`
	ISOWeek  string `json:"iso_week"`
}
