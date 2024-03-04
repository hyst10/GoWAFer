package controller

// LoginRequest 登录请求体
type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CaptchaID string `json:"captchaID" binding:"required"`
	Captcha   string `json:"captcha" binding:"required"`
}
