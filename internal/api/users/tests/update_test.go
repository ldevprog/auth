package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/ldevprog/auth/internal/api/users"
	"github.com/ldevprog/auth/internal/model"
	"github.com/ldevprog/auth/internal/service"
	serviceMocks "github.com/ldevprog/auth/internal/service/mocks"
	desc "github.com/ldevprog/auth/pkg/user_v1"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type usersServiceMockFunc func(mc *minimock.Controller) service.UsersService

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()

		serviceErr = fmt.Errorf("service error")

		req = &desc.UpdateRequest{
			Id: id,
			Name: &wrapperspb.StringValue{
				Value: name,
			},
			Email: &wrapperspb.StringValue{
				Value: email,
			},
		}

		userData = &model.UserChangable{
			Id:    id,
			Name:  &name,
			Email: &email,
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
				mock.UpdateMock.Expect(ctx, userData).Return(nil)
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
				mock.UpdateMock.Expect(ctx, userData).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			usersServiceMock := tt.usersServiceMock(mc)
			api := users.NewImplementation(usersServiceMock)

			res, err := api.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)
		})
	}
}
