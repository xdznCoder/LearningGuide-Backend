package model

type Post struct {
	BaseModel

	UserId     int32  `json:"user_id" gorm:"not null"`
	Category   string `json:"category" gorm:"type:varchar(50);not null"`
	Title      string `json:"title" gorm:"column:title;type:varchar(255);not null"`
	Content    string `json:"content" gorm:"type:longtext"`
	Desc       string `json:"desc" gorm:"type:text"`
	Image      string `json:"image" gorm:"type:varchar(255);not null"`
	LikeNum    int32  `json:"like_num" gorm:"not null;default=0"`
	FavNum     int32  `json:"fav_num" gorm:"not null;default=0"`
	CommentNum int32  `json:"comment_num" gorm:"not null;default=0"`
}
