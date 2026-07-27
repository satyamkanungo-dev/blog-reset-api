package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/Error"
)

var (
	AccessTokenExp  = time.Now().UTC().Add(15 * time.Minute)   //min
	RefreshTokenExp = time.Now().UTC().Add(7 * 24 * time.Hour) // days
)

func NewJWT(expTime time.Time, userId, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "blog-api",
		Subject:   userId,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(expTime),
	})

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// validate jwt
func ValidateJWT(tokenStr, secret string) (string, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", Error.ErrInvalidToken
	}

	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return "", Error.ErrTokenExpire
	}

	return claims.Subject, nil
}
