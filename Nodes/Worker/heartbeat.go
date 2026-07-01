package main

import (
	"context"
	pb "hadoop/gRPC/pb"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func SendHeartbeatToMaster(client pb.HealthServiceClient, node_id string) {
	req := &pb.HeartbeatRequest{
		WorkerId:  node_id,
		Ip:        GetOutboundIP().String(),
		Timestamp: time.Now().Unix(),
		Resources: &pb.NodeResources{
			TotalStorage: nodeManager.GetTotalStorage(),
			UsedStorage:  nodeManager.GetUsedStorage(),
			FreeStorage:  nodeManager.GetFreeStorage(),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.SendHeartbeat(ctx, req)
	if err != nil {
		log.Fatalf("Error sending heartbeat : %v", err)
	}
	log.Printf("Response from master %v", res.GetAcknowledge())
}

func HeartbeatLoop(conn *grpc.ClientConn, interval time.Duration, node_id string) {
	client := pb.NewHealthServiceClient(conn)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		SendHeartbeatToMaster(client, node_id)

	}

}

func startConnectionAndHeartbeatLoop(node_id *string) {
	const time_interval = 3 * time.Second
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("No connection found")
	}
	defer conn.Close()
	HeartbeatLoop(conn, time_interval, *node_id)
}
