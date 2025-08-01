package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ldevprog/auth/internal/model"
)

func GenerateToken(info *model.UserInfoForClaims, secretKey []byte, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		Username: info.Username,
		Role:     info.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
