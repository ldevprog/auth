package service

import (
	"context"

	"github.com/levon-dalakyan/auth/internal/model"
)

type UsersService interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.UserFullNoPass, error)
	Update(ctx context.Context, userData *model.UserChangable) error
	Delete(ctx context.Context, id int64) error
}
