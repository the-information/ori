package models

import (
	"time"
)

type AuthRequest struct {
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"length(6|999)"`
}

type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}
