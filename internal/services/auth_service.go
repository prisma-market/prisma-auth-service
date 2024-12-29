package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/config"
	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/models"
	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/repository/mongodb"
	"github.com/kihyun1998/prisma-market/prisma-auth-service/pkg/utils"
)

type AuthService struct {
	repo      *mongodb.AuthRepository
	jwtSecret string
	jwtExpiry int
	config    *config.Config // WebAppURL 등의 설정을 위해 필요
}

// NewAuthService AuthService 생성자
func NewAuthService(repo *mongodb.AuthRepository, config *config.Config) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: config.JWTSecret,
		jwtExpiry: config.JWTExpires,
		config:    config,
	}
}

// RegisterUser 사용자 등록
func (s *AuthService) RegisterUser(ctx context.Context, req *models.RegisterRequest) error {
	// 입력값 검증
	if err := validateRegisterRequest(req); err != nil {
		return err
	}

	// 비밀번호 해싱
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// 사용자 생성
	user := &models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	return s.repo.CreateUser(ctx, user)
}

// LoginUser 사용자 로그인
func (s *AuthService) LoginUser(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// 사용자 조회
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// 비밀번호 확인
	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 마지막 로그인 시간 업데이트
	if err := s.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		// 로깅만 하고 계속 진행
		// TODO: 로깅 추가
	}

	// JWT 토큰 생성
	token, err := utils.GenerateJWT(user.ID.Hex(), user.Email, s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:     token,
		ExpiresIn: s.jwtExpiry * 3600, // 시간을 초로 변환
	}, nil
}

func (s *AuthService) InitiatePasswordReset(ctx context.Context, req *models.ForgotPasswordRequest) error {
	// 사용자 조회
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return nil // 보안을 위해 사용자가 없어도 성공으로 처리
	}

	// 재설정 토큰 생성
	token := utils.GenerateRandomToken(32)
	expiry := time.Now().Add(1 * time.Hour)

	// 토큰 저장
	if err := s.repo.UpdateResetToken(ctx, user.Email, token, expiry); err != nil {
		return err
	}

	// 이메일 발송
	resetURL := fmt.Sprintf("%s/reset-password", s.config.WebAppURL)
	return s.emailService.SendPasswordResetEmail(user.Email, token, resetURL)
}

func (s *AuthService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	// 비밀번호 유효성 검사
	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		return err
	}

	// 비밀번호 해시화
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// 비밀번호 업데이트
	return s.repo.ResetPassword(ctx, req.Token, hashedPassword)
}

func (s *AuthService) SendVerificationEmail(ctx context.Context, email string) error {
	// 사용자 조회
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// 이미 인증된 경우
	if user.EmailVerified {
		return errors.New("email already verified")
	}

	// 인증 토큰 생성
	token := utils.GenerateRandomToken(32)
	expiry := time.Now().Add(24 * time.Hour)

	// 토큰 저장
	if err := s.repo.UpdateEmailVerificationToken(ctx, user.ID, token, expiry); err != nil {
		return err
	}

	// 이메일 발송
	verifyURL := fmt.Sprintf("%s/verify-email", s.config.WebAppURL)
	return s.emailService.SendVerificationEmail(user.Email, token, verifyURL)
}

func (s *AuthService) VerifyEmail(ctx context.Context, req *models.VerifyEmailRequest) error {
	return s.repo.VerifyEmail(ctx, req.Token)
}

// 입력값 검증 함수
func validateRegisterRequest(req *models.RegisterRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	// TODO: 이메일 형식 검증 추가
	return nil
}
