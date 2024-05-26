package types

type LoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required_if=Type 1 ,mobile"`
	Password  string `form:"password" json:"password" binding:"required_if=Type 1,omitempty,min=3,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,len=6"`
	CaptchaId string `form:"captchaId" json:"captchaId" binding:"required_if=Type 1"`
	Type      int    `form:"type" json:"type" binding:"required,oneof=1 2"` // 1表示密码登录, 2表示验证码登录
}

type RegisterForm struct {
	Username string `form:"username" json:"username" binding:"required,min=1,max=10"`
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=20"`
	Captcha  string `form:"captcha" json:"captcha" binding:"required,len=6"`
}

type UpdateUserForm struct {
	Username string `form:"username" json:"username" binding:"required,min=1,max=10"`
	Avatar   string `form:"avatar" json:"avatar" binding:"required,url"`
}

type UrlForm struct {
	Url string `form:"url" json:"url" binding:"required"`
}

type ResetPasswordForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=20"`
	Captcha  string `form:"captcha" json:"captcha" binding:"required,len=6"`
}

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Type   int    `form:"type" json:"type" binding:"required,oneof=1 2 3"` // 1表示注册, 2表示登录, 3表示修改密码
	Role   int    `form:"role" json:"role" binding:"required,oneof=1 2"`   // 1表示用户, 2表示管理员
}
