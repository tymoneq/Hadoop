package main

import (
	"context"

	pb "hadoop/Nodes/_proto/pb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SendHeartbeatToMaster(client pb.HealthServiceClient) {
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

func startHeartbeatLoop(conn *grpc.ClientConn, interval time.Duration) {
	client := pb.NewHealthServiceClient(conn)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			SendHeartbeatToMaster(client)
		}
	}

}

func main() {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("No connection found")
	}
	defer conn.Close()

	go startHeartbeatLoop(conn, 3*time.Second)
	log.Println("Worker started")
	select {}
}
