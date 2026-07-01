package main

import pb "hadoop/gRPC/pb"

type server struct {
	pb.UnimplementedHealthServiceServer
	pb.UnimplementedFileWritingMetadataServiceServer
}
