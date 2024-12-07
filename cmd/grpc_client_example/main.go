package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	desc "github.com/rkohnovets/go-auth/api/user_v1"
	utils "github.com/rkohnovets/go-auth/pkg/utils"
)

func main() {
	// grpc_server_address := "localhost:50051"
	grpc_server_address := "213.171.14.238:9090"
	conn, err := grpc.Dial(grpc_server_address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewUserV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var userId int64 = 1
	r, err := c.Get(ctx, &desc.IdRequest{Id: userId})
	if err != nil {
		log.Fatalf("failed to get user by id: %v", err)
	}

	result, err := utils.GetObjectJsonString(r)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	log.Println("value:\n" + string(result))
}
