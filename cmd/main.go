package main

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/brianvoe/gofakeit/v7"
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
	log.Printf("User id: %d", req.GetId())

	date := timestamppb.New(gofakeit.Date())
	return &desc.GetResponse{
		Id:        req.GetId(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      desc.Role_ADMIN,
		CreatedAt: date,
		UpdatedAt: date,
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("Updating user: %d, name: %s, email: %s", req.GetId(), req.GetName(), req.GetEmail())

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Deleting user: %d", req.GetId())

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
