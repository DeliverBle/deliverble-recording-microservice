# deliverble-recording-microservice
Recording S3 upload dedicated Microservice to deliverble restful server

## Architecture
![architecture.png](https://user-images.githubusercontent.com/41055141/203263203-4dfc5793-8472-40fd-8fb2-d1eef2c07b29.png)

## This module?
### does
* connect your mp3 file in byte format to communicate with deliverble s3 bucket
* only [deliverble restful server](https://github.com/DeliverBle/deliverble-backend-nestjs) is allowed to communicate with the MSA module
### does not
* provide any other functionality other than uploading mp3 file in s3 bucket
* delete feature or change naming strategy should be done in the restful server (nothing to do with bucket itself)

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

## Docker Commands
```
docker buildx build --platform linux/amd64 -f ./Dockerfile -t deliverble-recording-microservice .
```
