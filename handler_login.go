package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/Bention99/chirpy/internal/auth"
)

type validateParams struct {
	Email string `json:"email"`
	Password string `json:"password"`
	ExpiresInSeconds int `json:"expires_in_seconds"`
}

type RespondWithUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token	string `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := validateParams{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
    	return
    }

	expiresIn := time.Hour

	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expiresIn = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	u, err := cfg.db.GetUserByEMail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	match, err := auth.CheckPasswordHash(params.Password, u.HashedPassword)
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return 
	}

	jwtToken, err := auth.MakeJWT(u.ID, cfg.secret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return 
	}

	respondWithJSON(w, http.StatusOK, RespondWithUser{
		ID: u.ID,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.CreatedAt.Time,
		Email: u.Email,
		Token: jwtToken,
	})
}