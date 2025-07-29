package repository

import (
	"context"

	"github.com/ldevprog/auth/internal/model"
)

type UsersRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.UserFullNoPass, error)
	Update(ctx context.Context, data *model.UserChangable) error
	Delete(ctx context.Context, id int64) error
}
