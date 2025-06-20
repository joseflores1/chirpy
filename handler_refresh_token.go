package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/joseflores1/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	bearerToken, errGetToken := auth.GeatBearerToken(r.Header)
	if errGetToken != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't extract refresh token from header", errGetToken)
		return
	}

	refreshToken, err := cfg.db.GetTokenByID(r.Context(), bearerToken)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve refresh token from database", err)
		return
	}

	userID, errGetUser := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken.Token)
	if errGetUser != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve user's id using refresh token", errGetUser)
		return
	}

	expTime := time.Duration(1) * time.Hour

	jwtToken, errMakeJWT := auth.MakeJWT(userID, cfg.secretJWTKey, expTime)
	if errMakeJWT != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't produce JWT token", errMakeJWT)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: jwtToken,
	})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GeatBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't extract refresh token from header", err)
		return
	}

	errRevoke := cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if errRevoke != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", errRevoke)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
