.PHONY: all build proto clean

all: build

proto:
	protoc --go_out=. --go-grpc_out=. proto/container.proto

build: proto
	go build -o bin/containerhubd cmd/containerhubd/main.go
	go build -o bin/containerhub cmd/containerhub/main.go

clean:
	rm -rf bin/*
