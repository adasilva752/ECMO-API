package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/adasilva752/ECMO-API/user"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type myUserServer struct {
	user.UnimplementedUserServer
}

func (s *myUserServer) Create(ctx context.Context, req *user.CreateRequest) (*user.CreateResponse, error) {
	return &user.CreateResponse{
		ConfirmPass: "Here is the pass " + req.Password,
		Response:    "test",
	}, nil
}

func main() {
	// gRPC server
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}
	grpcServer := grpc.NewServer()
	service := &myUserServer{}
	user.RegisterUserServer(grpcServer, service)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// gRPC-Gateway server
	flag.Parse()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = user.RegisterUserHandlerFromEndpoint(ctx, mux, "localhost:8089", opts)
	if err != nil {
		log.Fatalf("failed to register gRPC-Gateway: %v", err)
	}

	log.Println("Starting HTTP/1.1 REST server on :8090")
	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
