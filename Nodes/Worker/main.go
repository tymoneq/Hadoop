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

	Manager := &NodeManager{}

	flag.Parse()
	log.Println("Worker started")

	go startConnectionAndHeartbeatLoop(node_id)
	go Manager.Initialize(node_id)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	log.Println("Usługi uruchomione w tle. Naciśnij Ctrl+C, aby wyłączyć.")
	<-stopChan
	log.Println("Otrzymano sygnał zamknięcia. Koniec pracy.")
}
