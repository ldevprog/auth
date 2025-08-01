package auth

import (
	"github.com/ldevprog/auth/internal/repository"
	"github.com/ldevprog/auth/internal/service"
)

const (
	refreshTokenSecretKeyEnvName     = "REFRESH_TOKEN_SECRET_KEY"
	refreshTokenExpirationMinEnvName = "REFRESH_TOKEN_EXPIRATION_MIN"
)

type serv struct {
	authRepository repository.AuthRepository
}

func NewService(authRepository repository.AuthRepository) service.AuthService {
	return &serv{
		authRepository: authRepository,
	}
}
