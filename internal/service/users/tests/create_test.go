package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/levon-dalakyan/auth/internal/model"
	"github.com/levon-dalakyan/auth/internal/repository"
	repoMocks "github.com/levon-dalakyan/auth/internal/repository/mocks"
	"github.com/levon-dalakyan/auth/internal/service/users"
	desc "github.com/levon-dalakyan/auth/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type usersRepositoryMockFunc func(mc *minimock.Controller) repository.UsersRepository

	type args struct {
		ctx  context.Context
		user *model.User
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Name()
		email     = gofakeit.Email()
		role      = desc.Role(gofakeit.RandomInt([]int{0, 1}))
		password  = gofakeit.Password(true, true, true, true, false, 10)
		password2 = gofakeit.Password(false, false, false, false, false, 15)

		repoErr = fmt.Errorf("repo error")

		user = &model.User{
			Name:            name,
			Email:           email,
			Role:            role,
			Password:        password,
			PasswordConfirm: password,
		}

		userPassNotMatch = &model.User{
			Name:            name,
			Email:           email,
			Role:            role,
			Password:        password,
			PasswordConfirm: password2,
		}
	)

	tests := []struct {
		name                string
		args                args
		want                int64
		err                 error
		usersRepositoryMock usersRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:  ctx,
				user: user,
			},
			want: id,
			err:  nil,
			usersRepositoryMock: func(mc *minimock.Controller) repository.UsersRepository {
				mock := repoMocks.NewUsersRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, user).Return(id, nil)
				return mock
			},
		},
		{
			name: "repository error case",
			args: args{
				ctx:  ctx,
				user: user,
			},
			want: 0,
			err:  repoErr,
			usersRepositoryMock: func(mc *minimock.Controller) repository.UsersRepository {
				mock := repoMocks.NewUsersRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, user).Return(0, repoErr)
				return mock
			},
		},
		{
			name: "passwords dont match err case",
			args: args{
				ctx:  ctx,
				user: userPassNotMatch,
			},
			want: 0,
			err:  errors.Errorf("passwords do not match"),
			usersRepositoryMock: func(mc *minimock.Controller) repository.UsersRepository {
				return repoMocks.NewUsersRepositoryMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			usersRepoMock := tt.usersRepositoryMock(mc)
			service := users.NewService(usersRepoMock)

			newID, err := service.Create(tt.args.ctx, tt.args.user)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.Equal(t, tt.err, err)
			}
			require.Equal(t, tt.want, newID)
		})
	}
}
