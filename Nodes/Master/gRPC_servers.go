package main

import (
	pb "hadoop/gRPC/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

const port string = ":50051"
const protocol string = "tcp"

type server struct {
	pb.UnimplementedHealthServiceServer
	pb.UnimplementedFileWritingMetadataServiceServer
}

func openConnection() {
	lis, err := net.Listen(protocol, port)
	if err != nil {
		log.Println("Something went wrong")
	}

	s := grpc.NewServer()

	pb.RegisterHealthServiceServer(s, &server{})
	pb.RegisterFileWritingMetadataServiceServer(s, &server{})

	log.Printf("Server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Server Error: %v", err)
	}

}
