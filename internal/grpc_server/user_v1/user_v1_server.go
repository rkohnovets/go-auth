package user_v1_server

import (
	"context"
	"log"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jackc/pgx/v5/pgxpool"
	dest "github.com/rkohnovets/go-auth/api/user_v1"
)

type User_v1_server struct {
	dest.UnimplementedUserV1Server
	DBPool *pgxpool.Pool
}

func (s *User_v1_server) Get(ctx context.Context, req *dest.IdRequest) (*dest.UserResponse, error) {
	userId := req.GetId()
	log.Printf("user id: %d", userId)

	return &dest.UserResponse{
		Id:        userId,
		Name:      gofakeit.BeerName(),
		Email:     gofakeit.Email(),
		Role:      dest.UserRoleEnum_USER,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

func (s *User_v1_server) Create(ctx context.Context, req *dest.UserRegisterRequest) (*dest.IdResponse, error) {
	userId := gofakeit.Int64()
	log.Printf("creating user, new id: %d", userId)

	return &dest.IdResponse{
		Id: userId,
	}, nil
}

func (s *User_v1_server) Delete(ctx context.Context, req *dest.IdRequest) (*emptypb.Empty, error) {
	userId := req.GetId()
	log.Printf("deleting user with id: %d", userId)

	return &emptypb.Empty{}, nil
}

func (s *User_v1_server) Update(ctx context.Context, req *dest.UserUpdateRequest) (*emptypb.Empty, error) {
	userId := req.GetId()
	log.Printf("updating user with id: %d", userId)

	return &emptypb.Empty{}, nil
}
