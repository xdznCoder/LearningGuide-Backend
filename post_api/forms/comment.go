package forms

type NewCommentForm struct {
	Content         string `json:"content" binding:"required"`
	UserId          int32  `json:"user_id" binding:"required"`
	PostId          int32  `json:"post_id" binding:"required"`
	ParentCommentId int32  `json:"parent_comment_id"`
}
