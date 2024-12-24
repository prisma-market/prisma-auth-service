package services

import (
	"context"
	"errors"

	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/models"
	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/repository/mongodb"
	"github.com/kihyun1998/prisma-market/prisma-auth-service/pkg/utils"
)

type AuthService struct {
	repo      *mongodb.AuthRepository
	jwtSecret string
	jwtExpiry int
}

// NewAuthService AuthService 생성자
func NewAuthService(repo *mongodb.AuthRepository, jwtSecret string, jwtExpiry int) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
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
