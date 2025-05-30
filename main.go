package main

import (
	"log"
	"net/http"
)

func main() {

	const port = "8080"

	// Define server struct and its handler
	serveMux := http.NewServeMux()
	serverStruct := &http.Server{
		Handler: serveMux,
		Addr: ":" + port,
	}	

	// Start the server
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(serverStruct.ListenAndServe())
}