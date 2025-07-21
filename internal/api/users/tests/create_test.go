package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/levon-dalakyan/auth/internal/api/users"
	"github.com/levon-dalakyan/auth/internal/model"
	"github.com/levon-dalakyan/auth/internal/service"
	serviceMocks "github.com/levon-dalakyan/auth/internal/service/mocks"
	desc "github.com/levon-dalakyan/auth/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type usersServiceMockFunc func(mc *minimock.Controller) service.UsersService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		role     = desc.Role(gofakeit.RandomInt([]int{0, 1}))
		password = gofakeit.Password(true, true, false, true, false, 10)

		serviceErr = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			Name:            name,
			Email:           email,
			Role:            role,
			Password:        password,
			PasswordConfirm: password,
		}

		user = &model.User{
			Name:            name,
			Email:           email,
			Role:            role,
			Password:        password,
			PasswordConfirm: password,
		}

		res = &desc.CreateResponse{
			Id: id,
		}
	)

	tests := []struct {
		name             string
		args             args
		want             *desc.CreateResponse
		err              error
		usersServiceMock usersServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			usersServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := serviceMocks.NewUsersServiceMock(mc)
				mock.CreateMock.Expect(ctx, user).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			usersServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := serviceMocks.NewUsersServiceMock(mc)
				mock.CreateMock.Expect(ctx, user).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			usersServiceMock := tt.usersServiceMock(mc)
			api := users.NewImplementation(usersServiceMock)

			createRes, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, createRes)
		})
	}
}
