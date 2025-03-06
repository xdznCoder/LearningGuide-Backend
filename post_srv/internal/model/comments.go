package model

type Comment struct {
	BaseModel

	UserId int32 `json:"user_id" gorm:"type:int;not null"`
	PostId int32 `json:"post_id" gorm:"type:int;not null"`

	ParentCommentId int32  `json:"parent_comment_id" gorm:"type:int;"`
	Content         string `json:"content" gorm:"type:text;not null"`
}

func (Comment) TableName() string {
	return "comment"
}
