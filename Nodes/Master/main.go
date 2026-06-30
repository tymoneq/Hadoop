package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	go openConnection()
	go checkHealth(5 * time.Second)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	log.Println("Usługi uruchomione w tle. Naciśnij Ctrl+C, aby wyłączyć.")
	<-stopChan
	log.Println("Otrzymano sygnał zamknięcia. Koniec pracy.")
}
