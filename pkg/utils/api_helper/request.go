package api_helper

// LoginRequest 登录请求体
type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CaptchaID string `json:"captchaId" binding:"required"`
	Captcha   string `json:"captcha" binding:"required"`
}

// CreateSqlInjectRequest 新增sql注入规则请求
type CreateSqlInjectRequest struct {
	Rule string `json:"rule" binding:"required"`
}

// CreateXssDetectRequest 新增xss攻击防御规则请求
type CreateXssDetectRequest struct {
	Rule string `json:"rule" binding:"required"`
}
