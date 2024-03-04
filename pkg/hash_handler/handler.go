package hash_handler

import "golang.org/x/crypto/bcrypt"

// EncryptPassword 加盐哈希加密明文密码
func EncryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// ValidatePassword 验证明文密码与哈希密码是否匹配
func ValidatePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
