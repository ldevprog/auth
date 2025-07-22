package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/levon-dalakyan/auth/internal/model"
	"github.com/levon-dalakyan/auth/internal/repository"
	repoMocks "github.com/levon-dalakyan/auth/internal/repository/mocks"
	"github.com/levon-dalakyan/auth/internal/service/users"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type usersRepositoryMockFunc func(mc *minimock.Controller) repository.UsersRepository

	type args struct {
		ctx      context.Context
		userData *model.UserChangable
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()

		repoErr = fmt.Errorf("repo error")

		userData = &model.UserChangable{
			Id:    id,
			Name:  &name,
			Email: &email,
		}
	)

	tests := []struct {
		name                string
		args                args
		err                 error
		usersRepositoryMock usersRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:      ctx,
				userData: userData,
			},
			err: nil,
			usersRepositoryMock: func(mc *minimock.Controller) repository.UsersRepository {
				mock := repoMocks.NewUsersRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, userData).Return(nil)
				return mock
			},
		},
		{
			name: "repository error case",
			args: args{
				ctx:      ctx,
				userData: userData,
			},
			err: repoErr,
			usersRepositoryMock: func(mc *minimock.Controller) repository.UsersRepository {
				mock := repoMocks.NewUsersRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, userData).Return(repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			usersRepoMock := tt.usersRepositoryMock(mc)
			service := users.NewService(usersRepoMock)

			err := service.Update(ctx, userData)
			require.Equal(t, tt.err, err)
		})
	}
}
