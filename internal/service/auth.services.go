package service

import (
	"time"

	"github.com/satyamkanungo-dev/blog-rest-api/internal/auth"
)

type AuthSsrvice struct {
	secret string
}

func NewAuthService(s string) *AuthSsrvice {
	return &AuthSsrvice{secret: s}
}

func (as *AuthSsrvice) GetToken(exp time.Time, userId string) (string, error) {
	return auth.NewJWT(exp, userId, as.secret)
}

func (as *AuthSsrvice) ValidateToken(tokenStr string) (string, error) {
	return auth.ValidateJWT(tokenStr, as.secret)
}
