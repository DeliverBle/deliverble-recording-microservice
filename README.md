# deliverble-recording-microservice
Recording S3 CRUD dedicated Microservice to deliverble restful server

## How to start
```
go get -u google.golang.org/grpc
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
```
protoc -I=. \
	    --go_out . --go_opt paths=source_relative \
	    --go-grpc_out . --go-grpc_opt paths=source_relative \
	    protos/v1/user/user.proto
```
