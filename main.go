package main

import (
	"context"
	"log"
	"net"

	"github.com/adasilva752/ECMO-API/user"
	"google.golang.org/grpc"
)

type myUserServer struct {
	user.UnimplementedUserServer
}

func (s myUserServer) Create(ctx context.Context, req *user.CreateRequest) (*user.CreateResponse, error) {
	return &user.CreateResponse{
		Response:    "Welcome " + req.Username,
		ConfirmPass: "here is the pass" + req.Password,
		Data:        []byte("test"),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}
	serverRegistrar := grpc.NewServer()
	service := &myUserServer{}
	user.RegisterUserServer(serverRegistrar, service)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
	}
}
