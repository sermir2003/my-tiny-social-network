export PATH="$PATH:$(go env GOPATH)/bin"
protoc \
    --proto_path ./proto/post \
    --experimental_allow_proto3_optional \
    --go_out=./proto/ \
    --go-grpc_out=./proto/ \
    post.proto
