package users

import (
	"context"

	"github.com/levon-dalakyan/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, userData *model.UserChangable) error {
	err := s.usersRepository.Update(ctx, userData)
	if err != nil {
		return err
	}

	return nil
}
