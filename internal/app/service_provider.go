package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	usersApi "github.com/levon-dalakyan/auth/internal/api/users"
	"github.com/levon-dalakyan/auth/internal/closer"
	"github.com/levon-dalakyan/auth/internal/config"
	"github.com/levon-dalakyan/auth/internal/repository"
	usersRepository "github.com/levon-dalakyan/auth/internal/repository/users"
	"github.com/levon-dalakyan/auth/internal/service"
	usersService "github.com/levon-dalakyan/auth/internal/service/users"
)

type serviceProvider struct {
	pgPool *pgxpool.Pool

	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	usersRepository repository.UsersRepository

	usersService service.UsersService

	usersImpl *usersApi.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) PGPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.Connect(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %s", err.Error())
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}

		closer.Add(func() error {
			pool.Close()
			return nil
		})

		s.pgPool = pool
	}

	return s.pgPool
}

func (s *serviceProvider) UsersRepository(ctx context.Context) repository.UsersRepository {
	if s.usersRepository == nil {
		s.usersRepository = usersRepository.NewRepository(s.PGPool(ctx))
	}

	return s.usersRepository
}

func (s *serviceProvider) UsersService(ctx context.Context) service.UsersService {
	if s.usersService == nil {
		s.usersService = usersService.NewService(s.UsersRepository(ctx))
	}

	return s.usersService
}

func (s *serviceProvider) UsersImpl(ctx context.Context) *usersApi.Implementation {
	if s.usersImpl == nil {
		s.usersImpl = usersApi.NewImplementation(s.UsersService(ctx))
	}

	return s.usersImpl
}
