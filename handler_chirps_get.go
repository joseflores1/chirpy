package main

import (
	"database/sql"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/joseflores1/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	var dbChirps []database.Chirp
	var errGetChirps error

	authorID := r.URL.Query().Get("author_id")
	if authorID != "" {
		parseID, errParse := uuid.Parse(authorID)
		if errParse != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id format", errParse)
			return
		}
		dbChirps, errGetChirps = cfg.db.GetChirpsByAuthor(r.Context(), parseID)
	} else {
		dbChirps, errGetChirps = cfg.db.GetChirps(r.Context())
	}

	if errGetChirps != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps from database", errGetChirps)
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

	sortParam := r.URL.Query().Get("sort")
	if sortParam == "desc" {
		sort.Slice(chirps, func(i, j int) bool {return chirps[i].CreatedAt.After(chirps[j].CreatedAt)})
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
		if errGetChirp == sql.ErrNoRows {
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
