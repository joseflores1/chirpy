package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {

	// Define consts
	const rootPath = "."
	const port = "8080"

	// Define apiConfig struct
	apiCfg := &apiConfig{
		fileServerHits: atomic.Int32{
		},
	}

	// Define mux
	serveMux := http.NewServeMux()
	// Define endpoints handlers
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(rootPath)))))
	serveMux.HandleFunc("/healthz", handlerReadiness)
	serveMux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("/reset", apiCfg.handlerResetHits)

	// Define server struct
	serverStruct := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// Start the server
	log.Printf("Serving files from %s on port: %s\n", rootPath, port)
	log.Fatal(serverStruct.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileServerHits.Load())))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(int32(1))
		next.ServeHTTP(w, r)
	})
}