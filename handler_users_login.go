package main

import (
	"encoding/json"
	"net/http"

	"github.com/joseflores1/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", errDecode)
		return
	}

	user, errGetUser := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if errGetUser != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", errGetUser)	
		return
	}

	errCheckPwd := auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if errCheckPwd != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", errCheckPwd)	
		return
	}

	type response struct {
		User
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}