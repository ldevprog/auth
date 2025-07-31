package app

import (
	"context"
	"log"

	"github.com/ldevprog/platform-common/pkg/closer"
	"github.com/ldevprog/platform-common/pkg/db"
	"github.com/ldevprog/platform-common/pkg/db/pg"

	authApi "github.com/ldevprog/auth/internal/api/auth"
	usersApi "github.com/ldevprog/auth/internal/api/users"
	"github.com/ldevprog/auth/internal/config"
	"github.com/ldevprog/auth/internal/repository"
	authRepository "github.com/ldevprog/auth/internal/repository/auth"
	usersRepository "github.com/ldevprog/auth/internal/repository/users"
	"github.com/ldevprog/auth/internal/service"
	authService "github.com/ldevprog/auth/internal/service/auth"
	usersService "github.com/ldevprog/auth/internal/service/users"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig

	dbClient        db.Client
	usersRepository repository.UsersRepository
	authRepository  repository.AuthRepository

	usersService service.UsersService
	authService  service.AuthService

	usersImpl *usersApi.Implementation
	authImpl  *authApi.Implementation
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

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
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

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) UsersService(ctx context.Context) service.UsersService {
	if s.usersService == nil {
		s.usersService = usersService.NewService(s.UsersRepository(ctx))
	}

	return s.usersService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.AuthRepository(ctx))
	}

	return s.authService
}

func (s *serviceProvider) UsersImpl(ctx context.Context) *usersApi.Implementation {
	if s.usersImpl == nil {
		s.usersImpl = usersApi.NewImplementation(s.UsersService(ctx))
	}

	return s.usersImpl
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *authApi.Implementation {
	if s.authImpl == nil {
		s.authImpl = authApi.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}
