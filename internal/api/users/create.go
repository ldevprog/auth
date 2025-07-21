package users

import (
	"context"

	"github.com/levon-dalakyan/auth/internal/converter"

	desc "github.com/levon-dalakyan/auth/pkg/user_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	userId, err := i.usersService.Create(ctx, converter.ToUserFromDesc(req))
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: userId,
	}, nil
}
