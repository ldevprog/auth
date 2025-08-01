package auth

import (
	"context"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ldevprog/auth/internal/model"
	"github.com/ldevprog/auth/internal/utils"
)

func (s *serv) Login(ctx context.Context, credentials *model.Credentials) (string, error) {
	creds, err := s.authRepository.Login(ctx, credentials.Username)
	if err != nil {
		return "", err
	}

	if credentials.Password != creds.Password {
		return "", status.Error(codes.InvalidArgument, "wrong password")
	}

	refreshTokenSecretKey := os.Getenv(refreshTokenSecretKeyEnvName)
	refreshTokenExpirationMin, err := strconv.Atoi(os.Getenv(refreshTokenExpirationMinEnvName))
	if err != nil {
		return "", status.Error(codes.Internal, "failed to generate token dueto unknown expiration")
	}
	expireDuration := time.Duration(refreshTokenExpirationMin) * time.Minute

	refreshToken, err := utils.GenerateToken(
		&model.UserInfoForClaims{
			Username: creds.Username,
			Role:     creds.Role.String(),
		},
		[]byte(refreshTokenSecretKey),
		expireDuration,
	)
	if err != nil {
		return "", status.Error(codes.Internal, "failed to generate token")
	}

	err = s.authRepository.SaveRefreshToken(ctx, &model.TokenWithCredentials{
		Token:     refreshToken,
		UserId:    creds.Id,
		ExpiresAt: time.Now().UTC().Add(expireDuration),
	})
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to save refresh token: %v", err)
	}

	return refreshToken, nil
}
