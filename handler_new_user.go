package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/Bention99/chirpy/internal/auth"
	"github.com/Bention99/chirpy/internal/database"
)

type createUserParams struct {
	Email string `json:"email"`
	Password string `json:"password"`
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

	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return 
	}

	u, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
        Email:          params.Email,
        HashedPassword: hashedPW,
    })
	if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
        return
    }

	respondWithJSON(w, http.StatusCreated, User{
		ID: u.ID,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
		Email: u.Email,
	})
}