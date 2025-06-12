package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handleUserCreation(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode request body", errDecode)
		return
	}

	user, errCreateUser := cfg.db.CreateUser(r.Context(), params.Email)
	if errCreateUser != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create database user", errCreateUser)
		return
	}

	type response struct {
		User
	}
	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
	}) 
}