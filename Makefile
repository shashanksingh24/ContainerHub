.PHONY: all build proto clean

PROTO_DIR=proto
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

GENERATE_PROTO:
       protoc \
               --go_out=. \
               --go_opt=paths=source_relative \
               --go-grpc_out=. \
               --go-grpc_opt=paths=source_relative \
               $(PROTO_FILES)

all: build

build: GENERATE_PROTO
	go build -o bin/containerhubd cmd/containerhubd/main.go
	go build -o bin/containerhub cmd/containerhub/main.go

run: build
       ./bin/containerhubd
       ./bin/containerhub

clean:
       rm -rf bin/ proto/*.pb.go
