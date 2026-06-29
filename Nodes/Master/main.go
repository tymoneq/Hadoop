package main

import (
	"context"
	pb "hadoop/Nodes/_proto/pb"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

const port string = ":50051"
const protocol string = "tcp"

type SafeNodes struct {
	mu    sync.RWMutex
	nodes map[string]int64
}

func (s *SafeNodes) Set(key string, value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodes[key] = value
}
func (s *SafeNodes) Get(key string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.nodes[key]
	return val, ok

}

type server struct {
	pb.UnimplementedHealthServiceServer
}

var nodes = &SafeNodes{nodes: make(map[string]int64)}

func (s *server) SendHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	log.Printf("Receive heartbeat from worker %s at %v", req.GetWorkerId(), req.GetTimestamp())
	nodes.Set(req.GetWorkerId(), req.GetTimestamp())
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
		for key, val := range nodes.nodes {
			if time.Now().Unix()-val >= 10 {
				log.Printf("Node %s is down", key)
			}
		}
	}

}

func main() {

	go openConnection()
	go checkHealth(5 * time.Second)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	log.Println("Usługi uruchomione w tle. Naciśnij Ctrl+C, aby wyłączyć.")
	<-stopChan
	log.Println("Otrzymano sygnał zamknięcia. Koniec pracy.")
}
