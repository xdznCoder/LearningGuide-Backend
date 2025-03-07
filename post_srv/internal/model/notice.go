package model

type Notice struct {
	BaseModel

	Type      int32 `json:"type"  gorm:"type:tinyint;not null"`
	UserId    int32 `json:"user_id" gorm:"type:int;not null"`
	OwnerId   int32 `json:"owner_id" gorm:"type:int;not null"`
	PostId    int32 `json:"post_id" gorm:"type:int;not null"`
	CommentId int32 `json:"comment_id" gorm:"type:int;not null"`
	IsRead    bool  `json:"is_read" gorm:"type:tinyint;not null"`
}

func (Notice) TableName() string {
	return "notice"
}

const (
	NoticeTypeLikeToPost = iota + 1
	NoticeTypeFavToPost
	NoticeTypeCommentToPost
	NoticeTypeCommentToComment
)
