package main

import (
	"encoding/json"
	"net/http"

	"github.com/joseflores1/chirpy/internal/auth"
	"github.com/joseflores1/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateCredentials(w http.ResponseWriter, r *http.Request) {

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

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", errDecode)
		return
	}

	hashedPwd, errHash := auth.HashPassword(params.Password)
	if errHash != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", errHash)
		return
	}

	user, errUpdate := cfg.db.UpdateCredentials(r.Context(), database.UpdateCredentialsParams{
		HashedPassword: hashedPwd,
		Email:          params.Email,
		ID:             userID,
	})
	if errUpdate != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update credentials", errUpdate)
		return
	}

	type response struct {
		User
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}
