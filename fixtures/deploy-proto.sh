#/bin/bash
export PATH="$PATH:$(go env GOPATH)/bin"

protoc \
    --proto_path ./proto \
    --experimental_allow_proto3_optional \
    --go_out=./main_service/core \
    --go_out=./ugc_service/core \
    --go-grpc_out=./main_service/core \
    --go-grpc_out=./ugc_service/core \
    post.proto

protoc \
    --proto_path ./proto \
    --experimental_allow_proto3_optional \
    --go_out=./main_service/core \
    --go_out=./stats_service/core \
    reaction.proto

protoc \
    --proto_path ./proto \
    --experimental_allow_proto3_optional \
    --go_out=./main_service/core \
    --go_out=./stats_service/core \
    --go-grpc_out=./main_service/core \
    --go-grpc_out=./stats_service/core \
    stats.proto
