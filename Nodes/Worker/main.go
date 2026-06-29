package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	pb "hadoop/Nodes/_proto/pb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SendHeartbeatToMaster(client pb.HealthServiceClient, node_id string) {
	req := &pb.HeartbeatRequest{
		WorkerId:  node_id,
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

func startHeartbeatLoop(conn *grpc.ClientConn, interval time.Duration, node_id string) {
	client := pb.NewHealthServiceClient(conn)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		SendHeartbeatToMaster(client, node_id)

	}

}

func main() {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("No connection found")
	}
	defer conn.Close()

	go startHeartbeatLoop(conn, 3*time.Second, "node-01")
	go startHeartbeatLoop(conn, 3*time.Second, "node-02")
	go startHeartbeatLoop(conn, 20*time.Second, "node-03")
	log.Println("Worker started")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	log.Println("Usługi uruchomione w tle. Naciśnij Ctrl+C, aby wyłączyć.")
	<-stopChan
	log.Println("Otrzymano sygnał zamknięcia. Koniec pracy.")
}
