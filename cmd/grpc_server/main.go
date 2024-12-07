package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jackc/pgx/v5/pgxpool"
	desc "github.com/rkohnovets/go-auth/api/user_v1"
	serv "github.com/rkohnovets/go-auth/internal/grpc_server/user_v1"
)

const grpcPort = 50051

func main() {
	// postgresql connection settings
	config, err := pgxpool.ParseConfig("postgres://username:password@localhost:5432/mydb")
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute
	// create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer pool.Close()

	// start tcp listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen at %d port: %v", grpcPort, err)
	}

	// start grpc server
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &serv.User_v1_server{
		DBPool: pool,
	})

	log.Printf("grpc server starting to listen at %v", listener.Addr())

	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
