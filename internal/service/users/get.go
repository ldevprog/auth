package users

import (
	"context"

	"github.com/ldevprog/auth/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.UserFullNoPass, error) {
	user, err := s.usersRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
