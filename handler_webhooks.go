package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/joseflores1/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {

	apiKey, errApiKey := auth.GetAPIKey(r.Header)
	if errApiKey != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get API key from header", errApiKey)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", errors.New("invalid api key"))
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", errDecode)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, errUpdate := cfg.db.UpgradeMembership(r.Context(), params.Data.UserID)
	if errUpdate != nil {
		if errUpdate == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "User doesn't exist", errUpdate)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", errUpdate)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
