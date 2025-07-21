package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/levon-dalakyan/auth/internal/api/users"
	"github.com/levon-dalakyan/auth/internal/service"
	serviceMocks "github.com/levon-dalakyan/auth/internal/service/mocks"
	desc "github.com/levon-dalakyan/auth/pkg/user_v1"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type usersServiceMockFunc func(mc *minimock.Controller) service.UsersService

	type args struct {
		ctx context.Context
		req *desc.DeleteRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		serviceErr = fmt.Errorf("service error")

		req = &desc.DeleteRequest{
			Id: id,
		}
	)

	tests := []struct {
		name             string
		args             args
		want             *emptypb.Empty
		err              error
		usersServiceMock usersServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &emptypb.Empty{},
			err:  nil,
			usersServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := serviceMocks.NewUsersServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &emptypb.Empty{},
			err:  serviceErr,
			usersServiceMock: func(mc *minimock.Controller) service.UsersService {
				mock := serviceMocks.NewUsersServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			usersServiceMock := tt.usersServiceMock(mc)
			api := users.NewImplementation(usersServiceMock)

			res, err := api.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
