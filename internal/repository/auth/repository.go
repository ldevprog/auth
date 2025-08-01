package auth

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ldevprog/auth/internal/model"
	"github.com/ldevprog/auth/internal/repository"
	"github.com/ldevprog/auth/internal/repository/auth/converter"
	modelRepo "github.com/ldevprog/auth/internal/repository/auth/model"
	"github.com/ldevprog/platform-common/pkg/db"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{db: db}
}

func (r *repo) Login(ctx context.Context, username string) (*model.CredentialsWithId, error) {
	builderSelect := sq.Select("id", "password", "role").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{
			"username": username,
		})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build SQL query: %v", err)
	}

	q := db.Query{
		Name:     "auth_repository.Login",
		QueryRaw: query,
	}

	credentials := modelRepo.CredentialsWithId{}
	err = r.db.DB().ScanOneContext(ctx, &credentials, q, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read user: %v", err)
	}

	return converter.ToCredentialsWithIdFromRepo(&credentials), nil
}

func (r *repo) SaveRefreshToken(ctx context.Context, tokenWithCreds *model.TokenWithCredentials) error {
	builderInsert := sq.Insert("refresh_tokens").
		PlaceholderFormat(sq.Dollar).
		Columns("token", "user_id", "expires_at").
		Values(
			tokenWithCreds.Token,
			tokenWithCreds.UserId,
			tokenWithCreds.ExpiresAt,
		)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to build SQL query: %v", err)
	}

	q := db.Query{
		Name:     "auth_repository.SaveRefreshToken",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to save refresh token: %v", err)
	}

	log.Printf("inserted %d rows", res.RowsAffected())

	return nil
}

func (r *repo) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	return "", nil
}
