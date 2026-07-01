package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "hadoop/gRPC/pb"

	"google.golang.org/grpc"
)

const port string = ":50051"
const protocol string = "tcp"

func updateWorker(ctx context.Context, req *pb.HeartbeatRequest) {

	workerID := req.GetWorkerId()

	nodeMaster.GetHeartbeats().Set(workerID, req.GetTimestamp())
	nodeMaster.GetNodeStatus().Set(workerID, true)
	nodeMaster.UpdateWorkerManager(workerID, req.Resources)
	nodeMaster.GetNodesIPManager().Set(workerID, req.GetIp())

}

func (s *server) SendHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	log.Printf("Receive heartbeat from worker %s at %v", req.GetWorkerId(), req.GetTimestamp())

	updateWorker(ctx, req)

	return &pb.HeartbeatResponse{
		Acknowledge: true,
	}, nil
}

func openConnection() {
	lis, err := net.Listen(protocol, port)
	if err != nil {
		log.Println("Something went wrong")
	}

	s := grpc.NewServer()

	pb.RegisterHealthServiceServer(s, &server{})

	log.Printf("Server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Server Error: %v", err)
	}

}

func checkHealth(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		for key, val := range nodeMaster.GetHeartbeats().nodes {
			if time.Now().Unix()-val >= 10 {
				nodeMaster.GetNodeStatus().Set(key, false)
				log.Printf("Node %s is down", key)
			}
		}
	}

}
