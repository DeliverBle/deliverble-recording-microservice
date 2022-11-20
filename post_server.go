package main

import (
	"deliverble-recording-msa/preprocess"
	postpb "deliverble-recording-msa/protos/v1/post"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	portNumber := "8081"

	lis, err := net.Listen("tcp", ":"+portNumber)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	postpb.RegisterPostServer(grpcServer, &preprocess.PostServer{})

	log.Printf("start gRPC post server on %s port", portNumber)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
