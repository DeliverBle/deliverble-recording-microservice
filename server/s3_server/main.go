package main

import (
	"deliverble-recording-msa/preprocess"
	recordingpb "deliverble-recording-msa/protos/v1/recording"
	"github.com/labstack/echo"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	portNumber := "8020"

	lis, err := net.Listen("tcp", ":"+portNumber)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	exit := make(chan bool)

	go func() {
		log.Println("start gRPC recording server on " + portNumber + " port")
		grpcServer := grpc.NewServer()
		recordingpb.RegisterRecordingTaskServer(grpcServer, &preprocess.S3Server{})
		if err := grpcServer.Serve(lis); err != nil {
			exit <- true
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	go func() {
		e := echo.New()
		e.POST("/upload/v2", preprocess.UploadRecordingHandlerV2)
		e.POST("/upload", preprocess.UploadRecordingHandler)
		err := e.Start(":8000")
		if err != nil {
			exit <- true
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	<-exit
}
