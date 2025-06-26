package models

import (
	"time"
)

type User struct {
	ID           int       `json:"id" postgres:"id"`
	Email        string    `json:"email" postgres:"email"`
	Username     string    `json:"username" postgres:"username"`
	PasswordHash string    `json:"-" postgres:"password_hash"`
	CreatedAt    time.Time `json:"created_at" postgres:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" postgres:"updated_at"`
}

type RefreshToken struct {
	ID        int       `json:"id" postgres:"id"`
	UserID    int       `json:"user_id" postgres:"user_id"`
	Token     string    `json:"token" postgres:"token"`
	ExpiresAt time.Time `json:"expires_at" postgres:"expires_at"`
	CreatedAt time.Time `json:"created_at" postgres:"created_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
