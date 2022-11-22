FROM golang:1.18.2-alpine AS builder
MAINTAINER "Sigrid Jin From SOPT 30th Deliverble"

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

RUN go build deliverble-recording-msa/server/s3_server

EXPOSE 8020
EXPOSE 8000

# run build binary
CMD ["./s3_server"]
