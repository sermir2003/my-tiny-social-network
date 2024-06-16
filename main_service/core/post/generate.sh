export PATH="$PATH:$(go env GOPATH)/bin"
protoc \
    --proto_path ./post \
    --experimental_allow_proto3_optional \
    --go_out=. \
    --go-grpc_out=. \
    post.proto
