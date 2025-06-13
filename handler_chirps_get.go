package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetAllUsers(w http.ResponseWriter, r *http.Request)  {

	dbChirps, errGetChirps := cfg.db.GetChirps(r.Context())
	if errGetChirps != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get all chirps in the database", errGetChirps)
		return
	}

	chirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		chirps[i] = Chirp{
			ID: dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body: dbChirp.Body,
			User_ID: dbChirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, chirps)
}