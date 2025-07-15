package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	sq "github.com/Masterminds/squirrel"
	desc "github.com/levon-dalakyan/auth/pkg/user_v1"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedUserV1Server
	db *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	pass := req.GetPassword()
	passConfirm := req.GetPasswordConfirm()
	if pass != passConfirm {
		return nil, status.Errorf(codes.InvalidArgument, "passwords do not match")
	}

	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("id", "name", "email", "role", "password").
		Values(randInt64Positive(), req.GetName(), req.GetEmail(), req.GetRole(), pass).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build SQL query: %v", err)
	}

	var userId int64
	err = s.db.QueryRow(ctx, query, args...).Scan(&userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert user: %v", err)
	}

	return &desc.CreateResponse{
		Id: userId,
	}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build SQL query: %v", err)
	}

	var id int64
	var name, email string
	var role desc.Role
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = s.db.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read user: %v", err)
	}

	var updatedAtProto *timestamppb.Timestamp
	if updatedAt.Valid {
		updatedAtProto = timestamppb.New(updatedAt.Time)
	} else {
		updatedAtProto = nil
	}

	return &desc.GetResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      role,
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: updatedAtProto,
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})
	if req.GetName() != nil {
		builderUpdate = builderUpdate.Set("name", req.GetName().Value)
	}
	if req.GetEmail() != nil {
		builderUpdate = builderUpdate.Set("email", req.GetEmail().Value)
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to build SQL query: %v", err)
	}

	res, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	builderDelete := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to build SQL query: %v", err)
	}

	res, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	log.Printf("deleted %d rows", res.RowsAffected())

	return &emptypb.Empty{}, nil
}

func randInt64Positive() int64 {
	var b [8]byte
	rand.Read(b[:])
	u := int64(binary.LittleEndian.Uint64(b[:]))
	return int64(u & 0x7FFFFFFFFFFFFFFF)
}

func getDSN() string {
	port := os.Getenv("PG_PORT")
	dbname := os.Getenv("PG_DATABASE_NAME")
	user := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWORD")

	return fmt.Sprintf(
		"host=localhost port=%s dbname=%s user=%s password=%s sslmode=disable",
		port,
		dbname,
		user,
		pass,
	)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file", err)
	}
}

func main() {
	ctx := context.Background()
	dbDSN := getDSN()

	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{db: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
