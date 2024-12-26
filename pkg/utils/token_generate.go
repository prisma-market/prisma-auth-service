package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

// GenerateRandomToken 지정된 바이트 수의 랜덤 토큰 생성
func GenerateRandomToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		// 에러메시지 개선 예약
		return ""
	}
	return hex.EncodeToString(b)
}

// GenerateBase64Token base64 인코딩된 랜덤 토큰 생성
func GenerateBase64Token(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		// 에러메시지 개선 예약
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
