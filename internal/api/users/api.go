package users

import (
	"github.com/levon-dalakyan/auth/internal/service"
	desc "github.com/levon-dalakyan/auth/pkg/user_v1"
)

type Implementation struct {
	desc.UnimplementedUserV1Server
	usersService service.UsersService
}

func NewImplementation(usersService service.UsersService) *Implementation {
	return &Implementation{
		usersService: usersService,
	}
}
