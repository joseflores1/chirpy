package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerResetHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(int32(0))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset at 0"))
}