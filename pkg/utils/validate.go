package utils

import (
	"fmt"
	"net/mail"
	"regexp"
	"unicode"
)

// ValidateEmail 이메일 주소 유효성 검사
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidatePassword 비밀번호 유효성 검사
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower {
		return fmt.Errorf("password must contain both uppercase and lowercase letters")
	}
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// ValidatePasswordResetToken 비밀번호 재설정 토큰 유효성 검사
func ValidatePasswordResetToken(token string) error {
	if len(token) < 32 {
		return fmt.Errorf("invalid reset token")
	}
	// 토큰이 16진수 문자열인지 확인
	if match, _ := regexp.MatchString("^[0-9a-fA-F]+$", token); !match {
		return fmt.Errorf("invalid reset token format")
	}
	return nil
}

// ValidateVerificationToken 이메일 인증 토큰 유효성 검사
func ValidateVerificationToken(token string) error {
	if len(token) < 32 {
		return fmt.Errorf("invalid verification token")
	}
	// 토큰이 16진수 문자열인지 확인
	if match, _ := regexp.MatchString("^[0-9a-fA-F]+$", token); !match {
		return fmt.Errorf("invalid verification token format")
	}
	return nil
}
