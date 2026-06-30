package main

import (
	"log"
	"os"
	"sync"
)

const TOTAL_STORAGE = 4096 * 10000
const FILE_PATH = "data/"

type NodeManager struct {
	totalStorage uint64
	usedStorage  uint64
	freeStorage  uint64
	ChunkManager *LocalChunkManager
}
type MyError struct{}

type LocalChunk struct {
	ChunkID     string
	FilePath    string
	Size        uint64
	IsCorrupted bool
}

type LocalChunkManager struct {
	mu     sync.RWMutex
	chunks map[string]LocalChunk
}

func (m *MyError) Error() string {
	return "Not enough free space"
}

func NewLocalChunkManager() *LocalChunkManager {
	return &LocalChunkManager{
		chunks: make(map[string]LocalChunk),
	}
}

func (node *NodeManager) InitializeStorage(node_id *string) (bool, error) {
	if err := os.MkdirAll(FILE_PATH+"data-for-"+*node_id, os.ModePerm); err != nil {
		log.Fatal(err)
		return false, err
	}
	log.Printf("Storage Created Successfully for node %s", *node_id)

	return true, nil

}

func (node *NodeManager) Initialize(node_id *string) {
	node.totalStorage = TOTAL_STORAGE
	node.freeStorage = TOTAL_STORAGE
	node.usedStorage = 0
	node.ChunkManager = NewLocalChunkManager()
	node.InitializeStorage(node_id)

}

func (node *NodeManager) GetFreeStorage() uint64 {
	return node.freeStorage
}
func (node *NodeManager) GetUsedStorage() uint64 {
	return node.usedStorage
}
func (node *NodeManager) GetTotalStorage() uint64 {
	return node.totalStorage
}
func (node *NodeManager) GetChunkManager() *LocalChunkManager {
	return node.ChunkManager
}

func (node *NodeManager) SaveFile(file_size uint64) (bool, error) {
	if file_size > node.GetFreeStorage() {
		return false, &MyError{}
	} else {
		log.Println("Saving file to do")
		return true, nil
	}

}
