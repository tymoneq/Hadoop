package main

import (
	"context"
	pb "hadoop/Nodes/_proto/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedHealthServiceServer
}

func (s *server) SendHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	log.Printf("Receive heartbeat from worker %s at %v", req.GetWorkerId(), req.GetTimestamp())

	return &pb.HeartbeatResponse{
		Acknowledge: true,
	}, nil
}

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Println("Something went wrong")
	}

	s := grpc.NewServer()

	pb.RegisterHealthServiceServer(s, &server{})

	log.Println("Server listening on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Server Error: %v", err)
	}

}
