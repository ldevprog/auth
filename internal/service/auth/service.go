package auth

import (
	"github.com/ldevprog/auth/internal/repository"
	"github.com/ldevprog/auth/internal/service"
)

type serv struct {
	authRepository repository.AuthRepository
}

func NewService(authRepository repository.AuthRepository) service.AuthService {
	return &serv{
		authRepository: authRepository,
	}
}
