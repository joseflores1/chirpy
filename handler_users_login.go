package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/joseflores1/chirpy/internal/auth"
	"github.com/joseflores1/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	expTime := time.Hour
	jwtToken, errMakeToken := auth.MakeJWT(user.ID, cfg.secretJWTKey, expTime)
	if errMakeToken != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", errMakeToken)
		return
	}

	refreshToken, _ := auth.MakeRefreshToken()
	_, errRefreshToken := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if errRefreshToken != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", errRefreshToken)
		return
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        jwtToken,
		RefreshToken: refreshToken,
	})
}
