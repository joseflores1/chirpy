package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/joseflores1/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	queries        *database.Queries
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, errOpen := sql.Open("postgres", dbURL)
	if errOpen != nil {
		log.Fatal("Couldn't open connection to database: ", errOpen)
	}
	dbQueries := database.New(db)

	// Define consts
	const rootPath = "."
	const port = "8080"

	// Define apiConfig struct
	apiCfg := &apiConfig{
		fileServerHits: atomic.Int32{},
		queries:        dbQueries,
	}

	// Define mux
	serveMux := http.NewServeMux()
	// Define endpoints handlers
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(rootPath)))))

	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerResetHits)

	// Define server struct
	serverStruct := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// Start the server
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(serverStruct.ListenAndServe())
}
