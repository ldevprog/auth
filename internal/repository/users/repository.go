package users

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/levon-dalakyan/auth/internal/client/db"
	"github.com/levon-dalakyan/auth/internal/helpers"
	"github.com/levon-dalakyan/auth/internal/model"
	"github.com/levon-dalakyan/auth/internal/repository"
	"github.com/levon-dalakyan/auth/internal/repository/users/converter"
	modelRepo "github.com/levon-dalakyan/auth/internal/repository/users/model"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UsersRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, user *model.User) (int64, error) {
	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "name", "email", "role", "password").
		Values(
			helpers.RandInt64Positive(),
			user.Name,
			user.Email,
			user.Role,
			user.Password,
		).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, status.Errorf(codes.Internal, "failed to build SQL query: %v", err)
	}

	q := db.Query{
		Name:     "users_repository.Create",
		QueryRaw: query,
	}

	var userId int64
	err = r.db.DB().ScanOneContext(ctx, &userId, q, args...)
	if err != nil {
		return 0, status.Errorf(codes.Internal, "failed to insert user: %v", err)
	}

	return userId, nil
}

func (r *repo) Get(ctx context.Context, userId int64) (*model.UserFullNoPass, error) {
	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": userId})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build SQL query: %v", err)
	}

	q := db.Query{
		Name:     "users_repository.Get",
		QueryRaw: query,
	}

	user := modelRepo.UserFullNoPass{}
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read user: %v", err)
	}

	return converter.ToUserFullNoPassFromRepo(&user), nil
}

func (r *repo) Update(ctx context.Context, user *model.UserChangable) error {
	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": user.Id})
	if user.Name != nil {
		builderUpdate = builderUpdate.Set("name", *user.Name)
	}
	if user.Email != nil {
		builderUpdate = builderUpdate.Set("email", *user.Email)
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to build SQL query: %v", err)
	}

	q := db.Query{
		Name:     "users_repository.Update",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builderDelete := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to build SQL query: %v", err)
	}

	q := db.Query{
		Name:     "users_repository.Delete",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	log.Printf("deleted %d rows", res.RowsAffected())

	return nil
}
