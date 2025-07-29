package users

import (
	"github.com/ldevprog/auth/internal/repository"
	"github.com/ldevprog/auth/internal/service"
)

type serv struct {
	usersRepository repository.UsersRepository
}

func NewService(usersRepository repository.UsersRepository) service.UsersService {
	return &serv{
		usersRepository: usersRepository,
	}
}
