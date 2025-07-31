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

const (
	refreshTokenSecretKeyEnvName     = "REFRESH_TOKEN_SECRET_KEY"
	refreshTokenExpirationMinEnvName = "REFRESH_TOKEN_EXPIRATION_MIN"
)

func (s *serv) Login(ctx context.Context, credentials *model.Credentials) (string, error) {
	creds, err := s.authRepository.Login(ctx, credentials.Name)
	if err != nil {
		return "", err
	}

	if credentials.Password != creds.Password {
		return "", status.Errorf(codes.InvalidArgument, "wrong password")
	}

	refreshTokenSecretKey := os.Getenv(refreshTokenSecretKeyEnvName)
	refreshTokenExpirationMin, err := strconv.Atoi(os.Getenv(refreshTokenExpirationMinEnvName))
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to generate token dueto unknown expiration")
	}

	refreshToken, err := utils.GenerateToken(
		&model.UserInfoForClaims{
			Name: creds.Name,
			Role: creds.Role.String(),
		},
		[]byte(refreshTokenSecretKey),
		time.Duration(refreshTokenExpirationMin)*time.Minute,
	)
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to generate token")
	}

	return refreshToken, nil
}
