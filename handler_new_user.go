package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/google/uuid"
	)

type createUserParams struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

    params := createUserParams{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
    	return
    }

	u, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't write parameters in table", err)
    	return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID: u.ID,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.CreatedAt.Time,
		Email: u.Email,
	})
}