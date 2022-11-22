#FROM golang:alpine AS builder
#MAINTAINER "Sigrid Jin from Deliverble"
#ENV GO111MODULE=on \
#    CGO_ENABLED=0 \
#    GOOS=linux \
#    GOARCH=amd64
#
#WORKDIR /build
#
## download the required Go dependencies
#COPY go.mod ./
#COPY go.sum ./
#RUN go mod download
#COPY . ./build
#
#RUN go build -o server/s3_server/main.go .
#
#EXPOSE 80
#EXPOSE 1432
# run build binary
#CMD ["./deliverble-recording-microservice"]

FROM golang:1.17-alpine AS builder

ENV GO111MODULE on

WORKDIR /app

COPY ./ ./build

RUN go mod download

RUN go build -o server/s3_server/main.go .

FROM gcr.io/distroless/base-debian10 AS runner

WORKDIR /build

EXPOSE 8020
EXPOSE 8000

# run build binary
CMD ["./deliverble-recording-microservice"]
