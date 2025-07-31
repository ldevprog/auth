package users

import (
	"context"

	"github.com/ldevprog/auth/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) Create(ctx context.Context, user *model.User) (int64, error) {
	if user.Password != user.PasswordConfirm {
		return 0, status.Errorf(codes.InvalidArgument, "passwords do not match")
	}

	userId, err := s.usersRepository.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
