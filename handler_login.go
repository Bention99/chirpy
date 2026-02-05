package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/Bention99/chirpy/internal/auth"
	"github.com/Bention99/chirpy/internal/database"
)

type validateParams struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type RespondWithUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token	string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed bool	`json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := validateParams{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
    	return
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

	jwtToken, err := auth.MakeJWT(u.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return 
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create Refresh Token", err)
		return 
	}

	expiresAt := time.Now().UTC().Add(60 * 24 * time.Hour)

	refreshTokenEntry, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: u.ID,
		ExpiresAt: expiresAt,
	})

	respondWithJSON(w, http.StatusOK, RespondWithUser{
		ID: u.ID,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.CreatedAt.Time,
		Email: u.Email,
		Token: jwtToken,
		RefreshToken: refreshTokenEntry.Token,
		IsChirpyRed: u.IsChirpyRed,
	})
}