package auth

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/ldevprog/auth/internal/model"
	"github.com/ldevprog/auth/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	refreshTokenSecretKey := os.Getenv(refreshTokenSecretKeyEnvName)
	claims, err := utils.VerifyToken(refreshToken, []byte(refreshTokenSecretKey))
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, err.Error())
	}

	refreshTokenExpirationMin, err := strconv.Atoi(os.Getenv(refreshTokenExpirationMinEnvName))
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to generate token dueto unknown expiration")
	}
	newRefreshToken, err := utils.GenerateToken(
		&model.UserInfoForClaims{
			Username: claims.Username,
			Role:     claims.Role,
		},
		[]byte(refreshTokenSecretKey),
		time.Duration(refreshTokenExpirationMin)*time.Minute,
	)
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to generate token")
	}

	return newRefreshToken, nil
}
