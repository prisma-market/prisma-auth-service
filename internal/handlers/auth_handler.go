package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/config"
	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/models"
	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/repository/mongodb"
	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// NewAuthHandler AuthHandler 생성자
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	repo, err := mongodb.NewAuthRepository(cfg.MongoURI)
	if err != nil {
		// 실제 운영환경에서는 에러 처리를 더 우아하게 해야 함
		panic(err)
	}

	authService := services.NewAuthService(repo, cfg.JWTSecret, cfg.JWTExpires)
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
