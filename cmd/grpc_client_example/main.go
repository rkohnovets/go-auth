package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"

	desc "github.com/rkohnovets/go-auth/api/user_v1"
	utils "github.com/rkohnovets/go-auth/pkg/utils"
)

func readLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	return input
}

func readInt64() (int64, error) {
	input := readLine()
	number, err := strconv.ParseInt(input, 10, 64)
	return number, err
}

func main() {
	grpc_server_address := "localhost:50051"
	// grpc_server_address := "213.171.14.238:9090"
	conn, err := grpc.Dial(
		grpc_server_address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewUserV1Client(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		fmt.Println("options: (0) exit, (1) create user, (2) get user by id, (3) update user, (4) delete user")
		fmt.Print("input: ")

		switch readLine() {
		case "0":
			return
		case "1":
			r, err := c.Create(ctx, &desc.UserRegisterRequest{
				Name:            gofakeit.BeerName(),
				Email:           gofakeit.Email(),
				Password:        "123",
				PasswordConfirm: "123",
				Role:            desc.UserRoleEnum_USER,
			})
			if err != nil {
				log.Fatalf("failed to create user: %v", err)
			}
			log.Printf("created user, id: %d", r.Id)
		case "2":
			fmt.Print("user id: ")
			id, err := readInt64()
			if err != nil {
				log.Printf("failed to read user id: %v", err)
				continue
			}

			r, err := c.Get(ctx, &desc.IdRequest{Id: id})
			if err != nil {
				log.Fatalf("failed to get user by id: %v", err)
			}

			result, err := utils.GetObjectJsonString(r)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			log.Println("got user:\n" + result)
		case "3":
			fmt.Print("user id: ")
			id, err := readInt64()
			if err != nil {
				log.Printf("failed to read user id: %v", err)
				continue
			}

			_, err = c.Update(ctx, &desc.UserUpdateRequest{
				Id:    id,
				Name:  wrapperspb.String(gofakeit.BeerName()),
				Email: wrapperspb.String(gofakeit.Email()),
				Role:  desc.UserRoleEnum_USER,
			})
			if err != nil {
				log.Fatalf("failed to update user: %v", err)
			}
			log.Printf("updated user")
		case "4":
			fmt.Print("user id: ")
			id, err := readInt64()
			if err != nil {
				log.Printf("failed to read user id: %v", err)
				continue
			}

			if _, err = c.Delete(ctx, &desc.IdRequest{Id: id}); err != nil {
				log.Fatalf("failed to update user: %v", err)
			}
			log.Printf("updated user")
		}
	}
}
