package captcha_handler

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

// GenerateCaptcha 生成验证码ID和图片
func GenerateCaptcha() (string, string, error) {
	// 创建验证码选项
	driver := base64Captcha.DriverDigit{Height: 80, Width: 240, Length: 6, MaxSkew: 0.7, DotCount: 80}
	captcha := base64Captcha.NewCaptcha(&driver, store)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		return "", "", fmt.Errorf("验证码生成失败：%v", err)
	}
	return id, b64s, nil
}

// VerifyCaptcha 检验验证码是否正确
func VerifyCaptcha(captchaId, captchaValue string) bool {
	return store.Verify(captchaId, captchaValue, true)
}
