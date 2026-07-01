package main

import (
	"context"
	"hadoop/gRPC/pb"
	"io"
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

type server struct {
	pb.UnimplementedFileSendingServiceServer
}

var client pb.FileWritingMetadataServiceClient
var sendingClient pb.FileSendingServiceClient

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

func SendFileMetadataToMaster(fileMetadata *FileToSend) (*pb.FileMetadataResponse, error) {
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
	return res, err
}

func streamFileToWorker(filePath string, nodes []*pb.NodesID) bool {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("%v", err)
		return false
	}
	defer file.Close()

	stream, err := sendingClient.SendChunk(context.Background())
	if err != nil {
		log.Fatalf("%v", err)
		return false
	}
	buffer := make([]byte, 64*1024)
	chunk_id := 1

	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			stream.Send(&pb.ChunkData{ChunkId: "2137", IsLast: true})
			break
		} else if err != nil {
			log.Fatalf("Something went wrong %v", err)
			break
		}

		err = stream.Send(&pb.ChunkData{
			ChunkId: string(chunk_id),
			Data:    buffer[:bytesRead],
			IsLast:  false,
		})
		chunk_id++
	}

	status, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error %v", err)
		return false
	}
	if status.GetSuccess() {
		log.Printf("File Send Successfully")
		return true
	}
	return false

}

func sendFile(filepath string) bool {
	fileMetadata := GetFileMetadata(filepath)
	status, err := SendFileMetadataToMaster(fileMetadata)
	if err != nil {
		return false
	}
	if streamFileToWorker(filepath, status.Nodes) {
		return true
	}

	return false
}
