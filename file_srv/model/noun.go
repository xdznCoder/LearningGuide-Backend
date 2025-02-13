package model

type Noun struct {
	BaseModel

	Name     string `json:"name" gorm:"type:varchar(20);not null"`
	Content  string `json:"content" gorm:"type:varchar(50);not null"`
	CourseId int32  `json:"course_id" gorm:"type:int;not null"`
}
