package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerResetUsers(w http.ResponseWriter, r *http.Request) {

	const allowedPlatform = "dev"

	if cfg.platform != allowedPlatform {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in local development environment"))
		return
	}

	cfg.fileServerHits.Store(0)
	errorDelete := cfg.db.Reset(r.Context())
	if errorDelete != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset the database: " + errorDelete.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset at 0 and database reset to initial state"))
}
