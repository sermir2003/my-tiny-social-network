export PATH="$PATH:$(go env GOPATH)/bin"
protoc \
    --proto_path ./reaction \
    --go_out=. \
    reaction.proto
