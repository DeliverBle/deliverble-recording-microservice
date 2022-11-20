package main

import (
	"deliverble-recording-msa/client"
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

	userCli := client.GetUserClient("localhost:8080")
	grpcServer := grpc.NewServer()
	postpb.RegisterPostServer(grpcServer, &preprocess.PostServer{
		UserClient: userCli,
	})

	log.Printf("start gRPC post server on %s port", portNumber)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
