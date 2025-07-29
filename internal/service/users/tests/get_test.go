package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/ldevprog/auth/internal/model"
	"github.com/ldevprog/auth/internal/repository"
	repoMocks "github.com/ldevprog/auth/internal/repository/mocks"
	"github.com/ldevprog/auth/internal/service/users"
	desc "github.com/ldevprog/auth/pkg/user_v1"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type usersRepositoryMockFunc func(mc *minimock.Controller) repository.UsersRepository

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Name()
		email     = gofakeit.Email()
		role      = desc.Role(gofakeit.RandomInt([]int{0, 1}))
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		repoErr = fmt.Errorf("repo error")

		user = &model.UserFullNoPass{
			Id:        id,
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: true,
			},
		}
	)

	tests := []struct {
		name                string
		args                args
		want                *model.UserFullNoPass
		err                 error
		usersRepositoryMock usersRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: user,
			err:  nil,
			usersRepositoryMock: func(mc *minimock.Controller) repository.UsersRepository {
				mock := repoMocks.NewUsersRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(user, nil)
				return mock
			},
		},
		{
			name: "repository error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  repoErr,
			usersRepositoryMock: func(mc *minimock.Controller) repository.UsersRepository {
				mock := repoMocks.NewUsersRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			usersRepoMock := tt.usersRepositoryMock(mc)
			service := users.NewService(usersRepoMock)

			user, err := service.Get(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, user)
		})
	}
}
