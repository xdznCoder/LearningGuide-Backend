package forms

type NewPostForm struct {
	UserId   int32  `json:"user_id" binding:"required"`
	Category string `json:"category" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Desc     string `json:"desc"`
	Image    string `json:"image"`
}

type UpdatePostForm struct {
	Content string `json:"content"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Image   string `json:"image"`
}
