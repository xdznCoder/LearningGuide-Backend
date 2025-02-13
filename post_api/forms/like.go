package forms

type NewLikeForm struct {
	UserId int32 `json:"user_id" binding:"required"`
	PostId int32 `json:"post_id" binding:"required"`
}

type DeleteLikeForm struct {
	UserId int32 `json:"user_id" binding:"required"`
	PostId int32 `json:"post_id" binding:"required"`
}
