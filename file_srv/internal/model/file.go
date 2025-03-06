package model

type File struct {
	BaseModel

	FileName string `gorm:"type:varchar(255);not null"`
	FileType string `gorm:"type:varchar(100);not null"`
	FileSize int64  `gorm:"type:bigint;not null"`
	OssUrl   string `gorm:"type:varchar(255);not null"`
	Desc     string `gorm:"type:text"`
	UserId   int32  `gorm:"not null"`
	CourseId int32  `gorm:"not null"`
}

func (File) TableName() string {
	return "file"
}
