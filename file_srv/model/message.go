package model

type Message struct {
	BaseModel

	Content   string `json:"content" gorm:"column:content;type:text;"`
	Speaker   string `json:"speaker" gorm:"column:speaker;type:varchar(20);"`
	SessionID int32  `json:"session_id" gorm:"column:session_id;type:int;not null"`
	Type      int    `json:"type" gorm:"type:int"`
}

type Session struct {
	BaseModel

	CourseId int32  `json:"course_id" gorm:"column:course_id;type:int;not null"`
	Uuid     string `json:"uuid" gorm:"column:uuid;type:varchar(64);unique;not null"`
}
