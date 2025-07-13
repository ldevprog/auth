package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/brianvoe/gofakeit"
	desc "github.com/levon-dalakyan/auth/pkg/user_v1"
)

const grpcPort = 50051

var users = map[int64]*desc.GetResponse{}

type server struct {
	desc.UnimplementedUserV1Server
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Creating user: %s, %s", req.GetName(), req.GetEmail())

	id := gofakeit.Int64()
	date := timestamppb.New(gofakeit.Date())
	users[id] = &desc.GetResponse{
		Id:        id,
		Name:      req.GetName(),
		Email:     req.GetEmail(),
		Role:      req.GetRole(),
		CreatedAt: date,
		UpdatedAt: date,
	}

	return &desc.CreateResponse{
		Id: id,
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

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
