package forms

type NewFavForm struct {
	UserId int32 `json:"user_id" binding:"required"`
	PostId int32 `json:"post_id" binding:"required"`
}

type DeleteFavForm struct {
	UserId int32 `json:"user_id" binding:"required"`
	PostId int32 `json:"post_id" binding:"required"`
}
