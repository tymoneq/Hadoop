package main

import (
	"context"

	pb "hadoop/Nodes/_proto/pb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("No connection found")
	}
	defer conn.Close()

	client := pb.NewHealthServiceClient(conn)

	req := &pb.HeartbeatRequest{
		WorkerId:  "worker-node-01",
		Timestamp: time.Now().Unix(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.SendHeartbeat(ctx, req)
	if err != nil {
		log.Fatalf("Error sending heartbeat : %v", err)
	}
	log.Printf("Response from master %v", res.GetAcknowledge())
}
