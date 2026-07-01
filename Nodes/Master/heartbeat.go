package main

import (
	"context"
	"log"
	"time"

	pb "hadoop/gRPC/pb"
)

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
