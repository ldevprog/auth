package users

import (
	"context"

	"github.com/ldevprog/auth/internal/converter"
	desc "github.com/ldevprog/auth/pkg/user_v1"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := i.usersService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converter.ToGetResponseFromService(user), nil
}
