package forms

type CreateSessionForm struct {
	CourseId int32 `json:"course_id" binding:"required,min=0"`
}

type SendMessageForm struct {
	SessionId int32  `json:"session_id" binding:"required,min=0"`
	Content   string `json:"content" binding:"required"`
	Type      int    `json:"type" binding:"required"`
}

type Message struct {
	Content string `json:"content"`
	Type    int    `json:"type"`
}
