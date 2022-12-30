FROM golang:1.18.2-alpine AS builder
MAINTAINER "Sigrid Jin From SOPT 30th Deliverble"

# install dependencies
RUN apk update && apk add --no-cache \
    alpine-sdk \
    ffmpeg

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# download the required Go dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./

# test only ./test/FFMpegConvert_test.go
RUN go test -v ./test/FFMpegConvert_test.go

RUN go build deliverble-recording-msa/server/s3_server

EXPOSE 8020
EXPOSE 8000

# run build binary
CMD ["./s3_server"]
