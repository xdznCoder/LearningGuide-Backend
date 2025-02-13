package forms

type PasswordLoginForm struct {
	Email     string `forms:"email" json:"email" binding:"required"`
	Password  string `forms:"password" json:"password" binding:"required,min=6,max=20"`
	Captcha   string `forms:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `forms:"captcha_id" json:"captcha_id" binding:"required"`
}

type EmailForm struct {
	Email string `forms:"email" json:"email" binding:"required"`
	Type  uint   `forms:"type" json:"type"`
}

type RegisterForm struct {
	Email      string `forms:"email" json:"email" binding:"required"`
	Nickname   string `forms:"nickname" json:"nickname" binding:"required,min=1,max=20"`
	Password   string `forms:"password" json:"password" binding:"required,min=6,max=20"`
	RePassword string `forms:"repassword" json:"repassword" binding:"required,min=6,max=20"`
	Code       string `forms:"code" json:"code" binding:"required,min=5,max=5"`
}

type UpdateUserForm struct {
	Nickname string `forms:"nickname" json:"nickname"`
	Birthday string `forms:"birthday" json:"birthday"`
	Gender   string `forms:"gender" json:"gender" binding:"required,oneof=male female"`
	Image    string `forms:"image" json:"image"`
	Desc     string `forms:"desc" json:"desc"`
}

type ChangePasswordForm struct {
	OldPassword string `forms:"old_password" json:"old_password" binding:"required,min=6,max=20"`
	Password    string `forms:"password" json:"password" binding:"required,min=6,max=20"`
	RePassword  string `forms:"repassword" json:"repassword" binding:"required,min=6,max=20"`
}

const (
	Register = iota
	Login
)
