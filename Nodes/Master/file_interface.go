package main

import (
	"context"
	pb "hadoop/gRPC/pb"
	"log"
)

const BLOCK_SIZE int64 = 8096
const REPLICATION_FACTOR uint8 = 3

type MyError struct{}

func (m *MyError) Error() string {
	return "Not enough free space"
}
func addPadding(fileSize int64) int64 {
	var extraBlock int64 = 0
	if fileSize%BLOCK_SIZE != 0 {
		extraBlock = 1
	}
	return (fileSize/BLOCK_SIZE + extraBlock) * BLOCK_SIZE
}

func checkIfEnoughFreeSpace(fileSize int64) bool {
	if nodeMaster.GetFreeStorage() >= fileSize*int64(REPLICATION_FACTOR) {
		return true
	}
	return false
}

func getNodes(fileSize int64) ([]*pb.NodesID, error) {
	blocks := fileSize / BLOCK_SIZE
	nodes := []*pb.NodesID{}

	const timeout int64 = 1000

	for i := int64(0); i < timeout && blocks > 0; i++ {
		for key, val := range nodeMaster.GetNodeStatus().nodes {
			if val {
				freeStorage, _ := nodeMaster.GetWorkerManager().Get(key)
				if freeStorage.freeStorage >= BLOCK_SIZE {

					blocks--
					ip, _ := nodeMaster.GetNodesIPManager().Get(key)
					newNode := pb.NodesID{
						NodeId: key,
						NodeIp: ip,
					}
					nodes = append(nodes, &newNode)
				}
			}
			if blocks == 0 {
				break
			}
		}

	}

	if blocks != 0 {
		return nodes, &MyError{}
	}

	return nodes, nil

}

func (s *server) SendFileMetadata(ctx context.Context, req *pb.FileMetadataRequest) (*pb.FileMetadataResponse, error) {
	log.Println("Receive file metadata for saving")
	fileSize := addPadding(req.GetTotalSize())
	canSaveFile := checkIfEnoughFreeSpace(fileSize)

	if !canSaveFile {
		return &pb.FileMetadataResponse{
			Status: false,
			Nodes:  []*pb.NodesID{},
		}, nil
	}
	nodesToSave, err := getNodes(fileSize)

	if err != nil {
		return &pb.FileMetadataResponse{
			Status: false,
			Nodes:  []*pb.NodesID{},
		}, nil
	}
	return &pb.FileMetadataResponse{
		Status: true,
		Nodes:  nodesToSave,
	}, nil

}
