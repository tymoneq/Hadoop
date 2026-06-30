package main

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
	nodeStatus    *SafeMap[string, bool]
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
