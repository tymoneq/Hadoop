package main

import (
	pb "hadoop/gRPC/pb"
	"log"
)

const BLOCK_SIZE int64 = 8096
const REPLICATION_FACTOR uint8 = 3

type WorkerManager struct {
	totalStorage int64
	usedStorage  int64
	freeStorage  int64
}

func (wm *WorkerManager) Subtract(w *WorkerManager) *WorkerManager {
	wm.freeStorage -= w.freeStorage
	wm.totalStorage -= w.totalStorage
	wm.usedStorage -= w.usedStorage
	return wm

}

type Metadata struct {
	fileToChunks  *SafeMap[string, []string] // key file name value list of chunks
	chunksToNodes *SafeMap[string, []string] // key chunk number value list of nodes with this chunk
}

type NodeMaster struct {
	totalStorage  int64
	freeStorage   int64
	usedStorage   int64
	heartbeats    *SafeMap[string, int64]
	workerManager *SafeMap[string, WorkerManager]
	nodeStatus    *SafeMap[string, bool] // true alive false down
	metadata      *Metadata
}

func NewMetadata() *Metadata {
	return &Metadata{
		fileToChunks:  NewSaveMap[string, []string](),
		chunksToNodes: NewSaveMap[string, []string](),
	}

}

func NewNodeMaster() *NodeMaster {
	return &NodeMaster{
		totalStorage:  0,
		freeStorage:   0,
		usedStorage:   0,
		heartbeats:    NewSaveMap[string, int64](),
		workerManager: NewSaveMap[string, WorkerManager](),
		nodeStatus:    NewSaveMap[string, bool](),
		metadata:      NewMetadata(),
	}

}

func (n *NodeMaster) GetTotalStorage() int64 {
	return n.totalStorage
}

func (n *NodeMaster) GetFreeStorage() int64 {
	return n.freeStorage
}

func (n *NodeMaster) GetUsedStorage() int64 {
	return n.usedStorage
}

func (n *NodeMaster) GetHeartbeats() *SafeMap[string, int64] {
	return n.heartbeats
}

func (n *NodeMaster) GetWorkerManager() *SafeMap[string, WorkerManager] {
	return n.workerManager
}

func (n *NodeMaster) GetNodeStatus() *SafeMap[string, bool] {
	return n.nodeStatus
}
func (n *NodeMaster) GetMetadata() *Metadata {
	return n.metadata
}

func (n *NodeMaster) UpdateStorageInfo(difference *WorkerManager) {
	n.totalStorage += difference.totalStorage
	n.freeStorage += difference.freeStorage
	n.usedStorage += difference.usedStorage

	log.Printf("\nSTORAGE UPDATE\n %d %d %d", n.totalStorage, n.freeStorage, n.usedStorage)

}

func (n *NodeMaster) UpdateWorkerManager(node_id string, workerResources *pb.NodeResources) {
	w := &WorkerManager{
		totalStorage: workerResources.TotalStorage,
		usedStorage:  workerResources.UsedStorage,
		freeStorage:  workerResources.FreeStorage,
	}
	prev_val, _ := n.GetWorkerManager().Get(node_id)
	if *w != prev_val {
		n.GetWorkerManager().Set(node_id, *w)
		diff := w.Subtract(&prev_val)
		n.UpdateStorageInfo(diff)

	}

}

var nodeMaster *NodeMaster = NewNodeMaster()
