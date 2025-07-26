package app

import (
	"context"
	"log"

	"github.com/levon-dalakyan/platform-common/pkg/closer"
	"github.com/levon-dalakyan/platform-common/pkg/db"
	"github.com/levon-dalakyan/platform-common/pkg/db/pg"

	usersApi "github.com/levon-dalakyan/auth/internal/api/users"
	"github.com/levon-dalakyan/auth/internal/config"
	"github.com/levon-dalakyan/auth/internal/repository"
	usersRepository "github.com/levon-dalakyan/auth/internal/repository/users"
	"github.com/levon-dalakyan/auth/internal/service"
	usersService "github.com/levon-dalakyan/auth/internal/service/users"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig
	httpConfig config.HTTPConfig

	dbClient        db.Client
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

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) UsersRepository(ctx context.Context) repository.UsersRepository {
	if s.usersRepository == nil {
		s.usersRepository = usersRepository.NewRepository(s.DBClient(ctx))
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
