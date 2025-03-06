package model

type Like struct {
	BaseModel

	UserId int32 `json:"user_id" gorm:"type:int;not null"`
	PostId int32 `json:"post_id" gorm:"type:int;not null"`
}

func (Like) TableName() string {
	return "like"
}
