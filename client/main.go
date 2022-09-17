package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	userpb "grpc-practice/user"
	"log"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := userpb.NewUserManagerClient(conn)

	log.Println("Request start!")

	req := &userpb.UserRequest{Id: 1}

	res, err := client.Get(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling User request: %v", err)
	}
	if res.Error {
		log.Fatalf("Bad UserRequest: %v", res.Message)
	}

	log.Printf("Response from UserManager Server: %v\n", res)
}
