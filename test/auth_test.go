package test

import (
	"GoWAFer/pkg/hash_handler"
	"GoWAFer/pkg/utils/jwt_handler"
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"testing"
	"time"
)

func TestAdminLogin(t *testing.T) {
	// 默认管理员登录名
	username := "admin"
	// 默认管理员登录密码
	password := "123456"
	// 检查是否存在此用户
	user, err := adminService.FindAdminByUsername(username)
	if err != nil {
		t.Error(err)
		return
	}
	// 检查密码是否匹配
	if !hash_handler.ValidatePassword(user.Password, password) {
		t.Error(err)
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
	token := jwt_handler.GenerateJwt(jwtClaims, "test")

	// 检查用户的refreshToken是否过期，过期就重新生成
	decodedRefreshToke := jwt_handler.VerifyUserRefreshToken(user.RefreshToken, "test")
	if decodedRefreshToke == nil {
		// 重新生成refreshToken
		jwtClaims = jwt.NewWithClaims(
			jwt.SigningMethodHS256, jwt.MapClaims{
				"id":  user.ID,
				"iat": time.Now().Unix(),
				"iss": os.Getenv("ENV"),
				"exp": time.Now().Add(7200 * time.Hour).Unix(),
			})
		user.RefreshToken = jwt_handler.GenerateJwt(jwtClaims, "test")
	}

	log.Printf("登录成功！\ntoken:%s\nrefreshToken:%s", token, user.RefreshToken)
}
