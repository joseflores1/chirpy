package main

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	dbChirps, errGetChirps := cfg.db.GetChirps(r.Context())
	if errGetChirps != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get all chirps in the database", errGetChirps)
		return
	}

	chirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		chirps[i] = Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			User_ID:   dbChirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	chirpID, errParse := uuid.Parse(r.PathValue("chirpID"))
	if errParse != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp's id", errParse)
		return
	}

	chirp, errGetChirp := cfg.db.GetChirp(r.Context(), chirpID)
	if errGetChirp != nil {
		if strings.Contains(errGetChirp.Error(), "no rows in result set") {
			respondWithError(w, http.StatusNotFound, "Chirp doesn't exist", errGetChirp)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp", errGetChirp)
			return
		}
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		User_ID:   chirp.UserID,
	})
}
