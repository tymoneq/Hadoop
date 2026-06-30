package main

import "hadoop/Nodes/_proto/pb"

const REPLICATION_FACTOR uint8 = 3

type WorkerManager struct {
	totalStorage uint64
	usedStorage  uint64
	freeStorage  uint64
}

type Metadata struct {
	fileToChunks  *SafeMap[string, []string] // key file name value list of chunks
	chunksToNodes *SafeMap[string, []string] // key chunk number value list of nodes with this chunk
}

type NodeMaster struct {
	totalStorage  uint64
	freeStorage   uint64
	usedStorage   uint64
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

func (n *NodeMaster) GetTotalStorage() uint64 {
	return n.totalStorage
}

func (n *NodeMaster) GetFreeStorage() uint64 {
	return n.freeStorage
}

func (n *NodeMaster) GetUsedStorage() uint64 {
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

func (n *NodeMaster) UpdateWorkerManager(node_id string, workerResources *pb.NodeResources) {
	w := WorkerManager{
		totalStorage: workerResources.TotalStorage,
		usedStorage:  workerResources.UsedStorage,
		freeStorage:  workerResources.FreeStorage,
	}

	n.GetWorkerManager().Set(node_id, w)

}

var nodeMaster *NodeMaster = NewNodeMaster()
