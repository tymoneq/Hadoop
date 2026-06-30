package main

import (
	"os"
	"os/signal"
	"syscall"

	"flag"
	"log"
)

func main() {

	node_id := flag.String("node-id", "node-00", "node id of cluster")

	flag.Parse()
	log.Println("Worker started")

	go startConnectionAndHeartbeatLoop(node_id)
	go nodeManager.InitializeStorage(node_id)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	log.Println("Running processes in the background.")
	<-stopChan
	log.Println("Termination Signal Receive.")
}
