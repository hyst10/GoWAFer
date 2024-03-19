package controller

import (
	"GoWAFer/internal/service"
	"GoWAFer/pkg/captcha_handler"
	"GoWAFer/pkg/config"
	"GoWAFer/pkg/hash_handler"
	"GoWAFer/pkg/utils/api_handler"
	"GoWAFer/pkg/utils/jwt_handler"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

// AdminController 管理员Controller层接口
type AdminController struct {
	adminService *service.AdminService
	conf         *config.Config
}

// NewAdminController 实例化管理员Controller层接口
func NewAdminController(adminService *service.AdminService, conf *config.Config) *AdminController {
	return &AdminController{
		adminService: adminService,
		conf:         conf,
	}
}

// DoLogin godoc
// @Summary 处理登录
// @Description dologin
// @Tags Auth
// @Accept json
// @Produce json
// @Param api_handler.LoginRequest body api_handler.LoginRequest true "请求体"
// @Success 200 {object} api_handler.Response
// @Router /waf/api/v1/auth/dologin [post]
func (c *AdminController) DoLogin(g *gin.Context) {
	var req api_handler.LoginRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		api_handler.ClientErrorHandler(g, 40001)
		return
	}
	// 检查验证码是否匹配
	if !captcha_handler.VerifyCaptcha(req.CaptchaID, req.Captcha) {
		api_handler.ClientErrorHandler(g, 40002)
		return
	}
	// 检查是否存在此用户
	user, err := c.adminService.FindAdminByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api_handler.ClientErrorHandler(g, 40003)
			return
		}
		api_handler.InternalErrorHandler(g, err)
		return
	}
	// 检查密码是否匹配
	if !hash_handler.ValidatePassword(user.Password, req.Password) {
		api_handler.ClientErrorHandler(g, 40003)
		return
	}

	// 生成token、refreshToken 存储到session中
	jwtClaims := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"iat":      time.Now().Unix(),
			"iss":      os.Getenv("ENV"),
			"exp":      time.Now().Add(10 * time.Minute).Unix(),
		})
	token := jwt_handler.GenerateJwt(jwtClaims, c.conf.Secret.JwtSecretKey)

	// 检查用户的refreshToken是否过期，过期就重新生成
	decodedRefreshToke := jwt_handler.VerifyUserRefreshToken(user.RefreshToken, c.conf.Secret.JwtSecretKey)
	if decodedRefreshToke == nil {
		// 重新生成refreshToken
		jwtClaims = jwt.NewWithClaims(
			jwt.SigningMethodHS256, jwt.MapClaims{
				"id":  user.ID,
				"iat": time.Now().Unix(),
				"iss": os.Getenv("ENV"),
				"exp": time.Now().Add(7200 * time.Hour).Unix(),
			})
		user.RefreshToken = jwt_handler.GenerateJwt(jwtClaims, c.conf.Secret.JwtSecretKey)
	}

	session := sessions.Default(g)
	session.Set("token", token)
	session.Set("refreshToken", user.RefreshToken)
	err = session.Save()
	if err != nil {
		api_handler.InternalErrorHandler(g, err)
		return
	}

	// 验证成功，记录登录日期、登录IP、refreshToken
	user.LastLoginDate = time.Now()
	user.LastLoginIP = g.ClientIP()
	if err = c.adminService.UpdateAdminInfo(user); err != nil {
		api_handler.InternalErrorHandler(g, err)
		return
	}

	// 登录成功，成功响应
	g.JSON(http.StatusOK, api_handler.Response{Status: 0, Message: "success"})
}

// GetCaptcha godoc
// @Summary 获取图片验证码
// @Description 获取图片验证码
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} api_handler.Response
// @Router /waf/api/v1/auth/getCaptcha [get]
func (c *AdminController) GetCaptcha(g *gin.Context) {
	captchaID, captchaBs6, err := captcha_handler.GenerateCaptcha()
	if err != nil {
		api_handler.InternalErrorHandler(g, err)
		return
	}
	g.JSON(http.StatusOK, api_handler.Response{Status: 0, Message: "success", Data: map[string]string{"id": captchaID, "bs6": captchaBs6}})
}
