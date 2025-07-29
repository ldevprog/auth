package users

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ldevprog/auth/internal/converter"
	desc "github.com/ldevprog/auth/pkg/user_v1"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := i.usersService.Update(ctx, converter.ToUserChangableFromDesc(req))
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}
