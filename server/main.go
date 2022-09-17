package main

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	userpb "grpc-practice/user"
	"io/ioutil"
	"log"
	"net"
	"strconv"
)

func (*server) Get(ctx context.Context, in *userpb.UserRequest) (*userpb.UserResponse, error) {
	log.Println("Request received!")
	id := in.GetId()

	// jsonの読み込み
	raw, err := ioutil.ReadFile("./users.json")
	if err != nil {
		log.Fatalf("Failed to open json file: %v", err)
	}

	var users Users
	if err = json.Unmarshal(raw, &users); err != nil {
		log.Fatalf("Failed to parse json: %v", err)
	}

	log.Println(users)

	errBool := false
	message := ""

	user, ok := users[strconv.Itoa(int(id))]
	if !ok {
		errBool = true
		message = "User not found"
	}

	return &userpb.UserResponse{
		Error:   errBool,
		Message: message,
		User: &userpb.User{
			Id:       user.Id,
			Name:     user.Name,
			Email:    user.Email,
			UserType: userpb.User_UserType(userpb.User_UserType_value[user.UserType]),
		},
	}, nil
}

type server struct {
	userpb.UnimplementedUserManagerServer
}

type User struct {
	Id       uint32 `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
}

type Users map[string]User

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	userpb.RegisterUserManagerServer(s, &server{})

	log.Println("Waiting for gRPC request ....")
	log.Println("-----------------------------")

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
