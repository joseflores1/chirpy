package main

import (
	"log"
	"net/http"
)

func main() {

	const rootPath = "."
	const port = "8080"

	// Define server structs and handlers
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(rootPath))))
	serveMux.HandleFunc("/healthz", handlerReadiness)
	serverStruct := &http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	// Start the server
	log.Printf("Serving files from %s on port: %s\n", rootPath, port)
	log.Fatal(serverStruct.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
