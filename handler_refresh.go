package main

import (
	"net/http"
	"time"
	"github.com/Bention99/chirpy/internal/auth"
)

type respondToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Bearer Token found", err)
		return
	}

	now := time.Now().UTC()

	tokenEntry, err := cfg.db.GetUserFromRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Bearer Token", err)
		return
	}

	if tokenEntry.RevokedAt.Valid || tokenEntry.ExpiresAt.Before(now) {
		respondWithError(w, http.StatusUnauthorized, "Invalid Bearer Token", nil)
		return
	}

	jwtToken, err := auth.MakeJWT(tokenEntry.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT", err)
		return 
	}

	respondWithJSON(w, http.StatusOK, respondToken{
		Token: jwtToken, 
	})
}