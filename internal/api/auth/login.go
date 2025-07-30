package auth

import (
	"context"

	"github.com/ldevprog/auth/internal/converter"
	desc "github.com/ldevprog/auth/pkg/auth_v1"
)

func (i *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	refreshToken, err := i.authService.Login(ctx, converter.ToCredentialsFromDesc(req))
	if err != nil {
		return nil, err
	}

	return &desc.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
