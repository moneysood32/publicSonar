package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// creates desired number of counter nodes at startup of server
	createCounters(3)

	coordinatorServer := &http.Server{
		Addr:    ":8080",
		Handler: CoordinatorHandler{},
	}
	fmt.Printf("Starting coordinator server...\n")
	log.Fatal(coordinatorServer.ListenAndServe())
}
