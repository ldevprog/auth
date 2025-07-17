package users

import (
	"context"

	"github.com/levon-dalakyan/auth/internal/converter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	desc "github.com/levon-dalakyan/auth/pkg/user_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, status.Errorf(codes.InvalidArgument, "passwords do not match")
	}

	userId, err := i.usersService.Create(ctx, converter.ToUserFromDesc(req))

	return &desc.CreateResponse{
		Id: userId,
	}, err
}
