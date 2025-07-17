package users

import (
	"github.com/levon-dalakyan/auth/internal/repository"
	"github.com/levon-dalakyan/auth/internal/service"
)

type serv struct {
	usersRepository repository.UsersRepository
}

func NewService(usersRepository repository.UsersRepository) service.UsersService {
	return &serv{
		usersRepository: usersRepository,
	}
}
