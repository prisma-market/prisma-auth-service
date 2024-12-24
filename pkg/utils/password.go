package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 비밀번호를 해시화
func HashPassword(password string) (string, error) {
	// bcrypt의 기본 cost 사용
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword 해시된 비밀번호와 일반 비밀번호 비교
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
