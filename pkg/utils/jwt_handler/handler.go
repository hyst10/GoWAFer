package jwt_handler

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"log"
)

// GenerateJwt 生成jwt
func GenerateJwt(claims *jwt.Token, secret string) (token string) {
	// 将 secret 转换成字节数组
	hmacSecretString := secret
	hmacSecret := []byte(hmacSecretString)

	// 通过 SignedString 方法生成签名后的JWT
	token, _ = claims.SignedString(hmacSecret)
	return
}

// VerifyUserToken 验证并解析用户token
func VerifyUserToken(token, secret string) *DecodeUserToken {
	// 将 secret 转换成字节数组
	hmacSecret := []byte(secret)

	// 解析JWT并验证签名
	decoded, err := jwt.Parse(
		token, func(token *jwt.Token) (interface{}, error) {
			return hmacSecret, nil
		})

	if err != nil {
		return nil
	}
	if !decoded.Valid {
		return nil
	}

	// 将解码后的JWT转换成结构体
	decodedClaims := decoded.Claims.(jwt.MapClaims)

	var decodedToken DecodeUserToken
	jsonString, _ := json.Marshal(decodedClaims)
	jsonErr := json.Unmarshal(jsonString, &decodedToken)
	if jsonErr != nil {
		log.Printf("用户token解析错误：%v", jsonErr)
	}
	return &decodedToken
}

// VerifyUserRefreshToken 验证并解析用户refreshToken
func VerifyUserRefreshToken(token, secret string) *DecodeUserRefreshToken {
	// 将 secret 转换成字节数组
	hmacSecret := []byte(secret)

	// 解析JWT并验证签名
	decoded, err := jwt.Parse(
		token, func(token *jwt.Token) (interface{}, error) {
			return hmacSecret, nil
		})

	if err != nil {
		return nil
	}
	if !decoded.Valid {
		return nil
	}

	// 将解码后的JWT转换成结构体
	decodedClaims := decoded.Claims.(jwt.MapClaims)

	var decodedToken DecodeUserRefreshToken
	jsonString, _ := json.Marshal(decodedClaims)
	jsonErr := json.Unmarshal(jsonString, &decodedToken)
	if jsonErr != nil {
		log.Printf("用户refreshToken解析错误：%v", jsonErr)
	}
	return &decodedToken
}
