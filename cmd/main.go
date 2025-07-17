package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/levon-dalakyan/auth/internal/config"
	"github.com/levon-dalakyan/auth/internal/converter"
	usersRepository "github.com/levon-dalakyan/auth/internal/repository/users"
	"github.com/levon-dalakyan/auth/internal/service"
	usersService "github.com/levon-dalakyan/auth/internal/service/users"
	desc "github.com/levon-dalakyan/auth/pkg/user_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedUserV1Server
	usersService service.UsersService
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, status.Errorf(codes.InvalidArgument, "passwords do not match")
	}

	userId, err := s.usersService.Create(ctx, converter.ToUserFromDesc(req))

	return &desc.CreateResponse{
		Id: userId,
	}, err
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := s.usersService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converter.ToGetResponseFromService(user), nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	err := s.usersService.Update(ctx, converter.ToUserChangableFromDesc(req))
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := s.usersService.Delete(ctx, req.GetId())
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	usersRepo := usersRepository.NewRepository(pool)
	usersServ := usersService.NewService(usersRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{usersService: usersServ})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
