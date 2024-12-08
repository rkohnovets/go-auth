package user_v1_server

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

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

	row := s.DBPool.QueryRow(
		context.Background(),
		"SELECT id, name, email, role, created_at, updated_at FROM users WHERE id = $1",
		userId,
	)

	var (
		updated, created time.Time
		result           dest.UserResponse
	)

	err := row.Scan(
		&result.Id,
		&result.Name,
		&result.Email,
		&result.Role,
		&created,
		&updated,
	)
	if err != nil {
		log.Fatalf("QueryRow failed: %v", err)
	}

	result.CreatedAt = timestamppb.New(created)
	result.UpdatedAt = timestamppb.New(updated)

	return &result, nil
}

func (s *User_v1_server) Create(ctx context.Context, req *dest.UserRegisterRequest) (*dest.IdResponse, error) {
	now := time.Now()
	row := s.DBPool.QueryRow(
		context.Background(),
		"INSERT INTO users (name, email, role, created_at, updated_at) "+
			"VALUES ($1, $2, $3, $4, $5) RETURNING id",
		req.Name, req.Email, req.Role, now, now,
	)

	var userId int64
	if err := row.Scan(&userId); err != nil {
		log.Fatalf("Failed to insert data: %v", err)
	}
	log.Println("Data inserted successfully")
	log.Printf("created user, new id: %d", userId)

	return &dest.IdResponse{
		Id: userId,
	}, nil
}

func (s *User_v1_server) Delete(ctx context.Context, req *dest.IdRequest) (*emptypb.Empty, error) {
	_, err := s.DBPool.Exec(context.Background(), "DELETE FROM users WHERE id = $1", req.GetId())
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
	log.Printf("deleted user with id: %d", req.Id)

	return &emptypb.Empty{}, nil
}

func (s *User_v1_server) Update(ctx context.Context, req *dest.UserUpdateRequest) (*emptypb.Empty, error) {
	userId := req.GetId()
	log.Printf("updating user with id: %d", userId)

	setParams := make([]any, 0)
	setString := ""
	if req.Email != nil {
		if len(setParams) != 0 {
			setString += ", "
		}
		setString += fmt.Sprintf("email = $%v", len(setParams)+1)
		setParams = append(setParams, req.Email.Value)
	}
	if req.Name != nil {
		if len(setParams) != 0 {
			setString += ", "
		}
		setString += fmt.Sprintf("name = $%v", len(setParams)+1)
		setParams = append(setParams, req.Name.Value)
	}

	setParams = append(setParams, userId)
	if len(setParams) != 0 {
		_, err := s.DBPool.Exec(
			context.Background(),
			"UPDATE users SET "+setString+" WHERE id = $"+strconv.Itoa(len(setParams)),
			setParams...,
		)
		if err != nil {
			log.Fatalf("Failed to update user: %v", err)
		}
		log.Printf("updated user with id: %d", req.Id)
	} else {
		log.Printf("not updated user with id %d (no fields to change)", req.Id)
	}

	req.GetEmail()

	return &emptypb.Empty{}, nil
}
