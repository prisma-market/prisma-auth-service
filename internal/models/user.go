package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email             string             `bson:"email" json:"email"`
	Password          string             `bson:"password" json:"-"`
	EmailVerified     bool               `bson:"email_verified" json:"email_verified"`
	EmailVerifyToken  string             `bson:"email_verify_token,omitempty" json:"-"`
	EmailVerifyExpiry time.Time          `bson:"email_verify_expiry,omitempty" json:"-"`
	ResetToken        string             `bson:"reset_token,omitempty" json:"-"`
	ResetTokenExpiry  time.Time          `bson:"reset_token_expiry,omitempty" json:"-"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
	Status            string             `bson:"status" json:"status"`
	LastLogin         *time.Time         `bson:"last_login,omitempty" json:"last_login,omitempty"`
}

// API 요청/응답 구조체
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

type ResendVerificationRequest struct {
	Email string `json:"email"`
}
