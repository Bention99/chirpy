package main

import (
	"encoding/json"
	"net/http"
	"github.com/Bention99/chirpy/internal/auth"
)

type validateParams struct {
	Email string `json:"email"`
	Password string `json:"password"`
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
	respondWithJSON(w, http.StatusOK, User{
		ID: u.ID,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.CreatedAt.Time,
		Email: u.Email,
	})
}