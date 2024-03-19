package middleware

import (
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/utils/api_handler"
	"GoWAFer/pkg/utils/jwt_handler"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func WafAPIAuthMiddleware(jwtSecretKey string, r *repository.AdminRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		token := session.Get("token")
		refreshToken := session.Get("refreshToken")

		if token != nil {
			tokenStr, _ := token.(string)
			decoded := jwt_handler.VerifyUserToken(tokenStr, jwtSecretKey)
			if decoded == nil {
				// token失效，检查refreshToken
				refreshTokenStr, _ := refreshToken.(string)
				decodedRefreshToken := jwt_handler.VerifyUserRefreshToken(refreshTokenStr, jwtSecretKey)

				if decodedRefreshToken == nil {
					// refreshToken也失效了，返回重新授权信息
					api_handler.ForbiddenHandler(c, "token已失效，请重新登录授权！")
					c.Abort()
					return
				}

				// refreshToken有效，重新生成token
				current, _ := r.FindByID(decodedRefreshToken.ID)
				jwtClaims := jwt.NewWithClaims(
					jwt.SigningMethodHS256, jwt.MapClaims{
						"username": current.Username,
						"iat":      time.Now().Unix(),
						"iss":      os.Getenv("ENV"),
						"exp":      time.Now().Add(10 * time.Minute).Unix(),
					})
				token = jwt_handler.GenerateJwt(jwtClaims, jwtSecretKey)
				session.Set("token", token)
				session.Save()
				c.Next()
				return
			}
			c.Next()
			return
		}
		api_handler.ForbiddenHandler(c, "您无权限！请先登录授权！")
		c.Abort()
		return
	}
}

func WafWebAuthMiddleware(jwtSecretKey string, r *repository.AdminRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		token := session.Get("token")
		refreshToken := session.Get("refreshToken")

		if token != nil {
			tokenStr, _ := token.(string)
			decoded := jwt_handler.VerifyUserToken(tokenStr, jwtSecretKey)
			if decoded == nil {
				// token失效，检查refreshToken
				refreshTokenStr, _ := refreshToken.(string)
				decodedRefreshToken := jwt_handler.VerifyUserRefreshToken(refreshTokenStr, jwtSecretKey)

				if decodedRefreshToken == nil {
					// refreshToken也失效了，返回重新授权信息
					c.Redirect(http.StatusFound, "/waf/login")
					c.Abort()
					return
				}

				// refreshToken有效，重新生成token
				current, _ := r.FindByID(decodedRefreshToken.ID)
				jwtClaims := jwt.NewWithClaims(
					jwt.SigningMethodHS256, jwt.MapClaims{
						"username": current.Username,
						"iat":      time.Now().Unix(),
						"iss":      os.Getenv("ENV"),
						"exp":      time.Now().Add(10 * time.Minute).Unix(),
					})
				token = jwt_handler.GenerateJwt(jwtClaims, jwtSecretKey)
				session.Set("token", token)
				session.Save()
				c.Next()
				return
			}
			c.Next()
			return
		}
		c.Redirect(http.StatusFound, "/waf/login")
		c.Abort()
		return
	}
}
