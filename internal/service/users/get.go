package users

import (
	"context"

	"github.com/levon-dalakyan/auth/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.UserFullNoPass, error) {
	userId, err := s.usersRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return userId, nil
}
