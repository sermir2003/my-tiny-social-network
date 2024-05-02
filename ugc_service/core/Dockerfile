FROM golang:1.22.2 AS builder
WORKDIR /app

RUN apt update && apt install --assume-yes tree

COPY . .
RUN mkdir -p ./proto/post
COPY ../../proto/post/post.proto ./proto/post

RUN tree .

RUN go mod download

RUN apt update && apt install --assume-yes protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN ./proto/generate.sh

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -trimpath -mod=readonly -o /go/bin/app

ENTRYPOINT [ "/go/bin/app" ]