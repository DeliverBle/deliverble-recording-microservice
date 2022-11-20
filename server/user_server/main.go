package main

import (
	"deliverble-recording-msa/preprocess"
	userpb "deliverble-recording-msa/protos/v1/user"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	portNumber := "8080"

	lis, err := net.Listen("tcp", ":"+portNumber)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServer(grpcServer, &preprocess.UserServer{})

	log.Printf("start gRPC user server on %s port", portNumber)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
