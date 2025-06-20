package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/joseflores1/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	bearerToken, errGetToken := auth.GeatBearerToken(r.Header)
	if errGetToken != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get JWT from header", errGetToken)
		return
	}

	userID, errValidateJWT := auth.ValidateJWT(bearerToken, cfg.secretJWTKey)
	if errValidateJWT != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", errValidateJWT)
		return
	}

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

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Can't delete another user's chirp", errors.New("forbidden chirp deletion request"))
		return
	}

	errDeleteChirp := cfg.db.DeleteChirp(r.Context(), chirpID)
	if errDeleteChirp != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", errDeleteChirp)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
