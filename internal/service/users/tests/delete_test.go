package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/ldevprog/auth/internal/repository"
	repoMocks "github.com/ldevprog/auth/internal/repository/mocks"
	"github.com/ldevprog/auth/internal/service/users"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type usersRepositoryMockFunc func(mc *minimock.Controller) repository.UsersRepository

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		repoErr = fmt.Errorf("repo error")
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
				ctx: ctx,
				id:  id,
			},
			err: nil,
			usersRepositoryMock: func(mc *minimock.Controller) repository.UsersRepository {
				mock := repoMocks.NewUsersRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "repository error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: repoErr,
			usersRepositoryMock: func(mc *minimock.Controller) repository.UsersRepository {
				mock := repoMocks.NewUsersRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			usersRepoMock := tt.usersRepositoryMock(mc)
			service := users.NewService(usersRepoMock)

			err := service.Delete(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.err, err)
		})
	}
}
