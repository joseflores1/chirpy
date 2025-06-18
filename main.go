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
	db             *database.Queries
	platform       string
	secretJWTKey   string
}

func main() {

	// Define consts
	const rootPath = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	dbConn, errOpen := sql.Open("postgres", dbURL)
	if errOpen != nil {
		log.Fatal("Couldn't open connection to database: ", errOpen)
	}
	dbQueries := database.New(dbConn)

	platformStr := os.Getenv("PLATFORM")
	if platformStr == "" {
		log.Fatal("PLATFORM must be set")
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY must be set")
	}

	// Define apiConfig struct
	apiCfg := &apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platformStr,
		secretJWTKey:   secretKey,
	}

	// Define mux
	serveMux := http.NewServeMux()
	// Define endpoints handlers
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(rootPath)))))

	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpCreation)

	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)

	serveMux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)

	serveMux.HandleFunc("POST /api/users", apiCfg.handleUserCreation)

	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerResetUsers)

	// Define server struct
	serverStruct := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	// Start the server
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(serverStruct.ListenAndServe())
}
