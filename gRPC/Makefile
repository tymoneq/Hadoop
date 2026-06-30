GOPATH_BIN := $(shell go env GOPATH)/bin
export PATH := $(GOPATH_BIN):$(PATH)

proto: heartbeat.proto
	rm pb/*
	protoc --go_out=. --go-grpc_out=. heartbeat.proto