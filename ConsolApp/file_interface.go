package main

import (
	"context"
	"hadoop/gRPC/pb"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FileToSend struct {
	totalSize int64
	fileName  string
}

var client pb.FileWritingMetadataServiceClient

func startConnection() *grpc.ClientConn {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("No connection found")
	}

	client = pb.NewFileWritingMetadataServiceClient(conn)
	return conn
}

func IsPathValid(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

func GetFileMetadata(filepath string) *FileToSend {
	splitStr := strings.Split("/", filepath)
	fileName := splitStr[len(splitStr)-1]
	info, _ := os.Stat(filepath)
	var fileSize int64 = info.Size()
	return &FileToSend{
		totalSize: fileSize,
		fileName:  fileName,
	}

}

func SendFileMetadataToMaster(fileMetadata *FileToSend) {
	req := &pb.FileMetadataRequest{
		TotalSize: fileMetadata.totalSize,
		FileName:  fileMetadata.fileName,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.SendFileMetadata(ctx, req)
	if err != nil {
		log.Fatalf("Error sending file metadata : %v\n", err)
	}

	log.Printf("Response from master %v\n", res)

}

func sendFile(filepath string) bool {
	fileMetadata := GetFileMetadata(filepath)
	SendFileMetadataToMaster(fileMetadata)

	return true
}
