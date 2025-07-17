package users

import (
	"context"

	"github.com/levon-dalakyan/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, user *model.User) (int64, error) {
	userId, err := s.usersRepository.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
