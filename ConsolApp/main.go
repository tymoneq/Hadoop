package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const BLOCK_SIZE int = 8096
const REPLICATION_FACTOR uint8 = 3

func main() {

	stopChan := make(chan os.Signal, 1)

	inputChan := make(chan string)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {

			if !scanner.Scan() {
				break
			}

			inputChan <- strings.TrimSpace(scanner.Text())
		}
	}()
	fmt.Println("Please enter absolute path to a file")

	for {
		select {
		case path := <-inputChan:
			fmt.Printf("Validating: %s\n", path)
			if IsPathValid(path) {
				fmt.Println("PATH VALID connecting to the master node")
				if sendFile(path) {
					fmt.Println("File successfully sent")
				} else {
					fmt.Println("There was an error")
				}

			} else {
				fmt.Printf("%s path is not valid\n", path)
			}
			fmt.Println("Please enter absolute path to a file")

		case <-stopChan:
			log.Println("Termination Signal Receive.")
			return
		}
	}
}
