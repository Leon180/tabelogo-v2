#!/bin/bash

# Install plugins if needed (uncomment if you want to automate installation)
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Ensure output directory exists
mkdir -p api/gen

# Generate Auth Service Protos
protoc --go_out=. --go_opt=module=github.com/Leon180/tabelogo-v2 \
    --go-grpc_out=. --go-grpc_opt=module=github.com/Leon180/tabelogo-v2 \
    api/proto/auth/v1/auth.proto

echo "Generated Auth Service Protos"

# Generate Spider Service Protos
protoc --go_out=. --go_opt=module=github.com/Leon180/tabelogo-v2 \
    --go-grpc_out=. --go-grpc_opt=module=github.com/Leon180/tabelogo-v2 \
    api/proto/spider/v1/spider.proto

echo "Generated Spider Service Protos"
