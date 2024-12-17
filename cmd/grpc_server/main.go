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
	"github.com/rkohnovets/go-auth/internal/config"
	serv "github.com/rkohnovets/go-auth/internal/grpc_server/user_v1"
)

func main() {
	ctx := context.Background()
	logger := log.Default()
	config := config.GetConfig(logger)

	// create database connection pool
	pool := createPgxPool(ctx, &config)
	defer pool.Close()

	// create tcp listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Listen.Port))
	if err != nil {
		log.Fatalf("failed to listen at %d port: %v", config.Listen.Port, err)
	}

	// create grpc server
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

func createPgxPool(ctx context.Context, config *config.Config) *pgxpool.Pool {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DatabaseName,
	)

	// postgresql connection settings
	dbConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}
	dbConfig.MaxConns = 10
	dbConfig.MinConns = 2
	dbConfig.MaxConnLifetime = time.Hour
	dbConfig.MaxConnIdleTime = 30 * time.Minute
	dbConfig.HealthCheckPeriod = time.Minute

	// create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	return pool
}
