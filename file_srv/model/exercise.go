package model

type Exercise struct {
	BaseModel

	CourseId int32  `json:"course_id" gorm:"column:course_id;type:int;not null"`
	Question string `json:"question" gorm:"column:question;type:text"`
	SectionA string `json:"sectionA" gorm:"column:section_a;type:text"`
	SectionB string `json:"sectionB" gorm:"column:section_b;type:text"`
	SectionC string `json:"sectionC" gorm:"column:section_c;type:text"`
	SectionD string `json:"sectionD" gorm:"column:section_d;type:text"`
	Answer   string `json:"answer" gorm:"column:answer;type:text"`
	Reason   string `json:"reason" gorm:"column:reason;type:text"`
	IsRight  string `json:"is_right" gorm:"column:is_right;type:varchar(20)"`
}
