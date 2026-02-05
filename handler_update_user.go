package main

import (
	"encoding/json"
	"net/http"
	"github.com/Bention99/chirpy/internal/auth"
	"github.com/Bention99/chirpy/internal/database"
)

type updateParams struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Bearer Token found", err)
		return
	}

	uID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect bearer Token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)

    params := updateParams{}
    err = decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
    	return
    }

	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return 
	}

	u, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
        Email:          params.Email,
        HashedPassword: hashedPW,
		ID:				uID,
    })
	if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
        return
    }

	respondWithJSON(w, http.StatusOK, User{
		ID: u.ID,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
		Email: u.Email,
		IsChirpyRed: u.IsChirpyRed,
	})
}