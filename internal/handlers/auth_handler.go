package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/config"
	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/models"
	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/repository/mongodb"
	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/services"
	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/services/email"
)

type AuthHandler struct {
	authService *services.AuthService
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// NewAuthHandler AuthHandler 생성자
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	// MongoDB Repository 초기화
	repo, err := mongodb.NewAuthRepository(cfg.MongoURI)
	if err != nil {
		// 에러처리 예약
		panic(err)
	}

	// Email Service 초기화
	emailService := email.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.SMTPFrom,
	)

	// Auth Service 초기화
	authService := services.NewAuthService(repo, emailService, cfg)

	return &AuthHandler{
		authService: authService,
	}
}

// Register 회원가입 핸들러
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.authService.RegisterUser(r.Context(), &req); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}

// Login 로그인 핸들러
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.authService.LoginUser(r.Context(), &req)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// sendError 에러 응답 전송 헬퍼 함수
func (h *AuthHandler) sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// ForgotPassword 비밀번호 재설정 요청
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.authService.InitiatePasswordReset(r.Context(), &req); err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If an account exists with that email, a password reset link has been sent",
	})
}

// ResetPassword 비밀번호 재설정
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.authService.ResetPassword(r.Context(), &req); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password has been reset successfully",
	})
}

// SendVerificationEmail 이메일 인증 메일 발송
func (h *AuthHandler) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	var req models.ResendVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.authService.SendVerificationEmail(r.Context(), req.Email); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Verification email has been sent",
	})
}

// VerifyEmail 이메일 인증 처리
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.authService.VerifyEmail(r.Context(), &req); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Email has been verified successfully",
	})
}
