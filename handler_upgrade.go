package main

import (
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"github.com/Bention99/chirpy/internal/auth"
)

type EventPayload struct {
	Event string    `json:"event"`
	Data  EventData `json:"data"`
}

type EventData struct {
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerUpgrade(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No API Key found", err)
		return
	}

	if apiKey != cfg.polkaKey {
			respondWithError(w, http.StatusUnauthorized, "Invalid API Key", err)
    	return
	}

	decoder := json.NewDecoder(r.Body)

    params := EventPayload{}
    err = decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
    	return
    }

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Wrong Event", err)
    	return
	}

	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
    	return
	}

	respondWithStatus(w, http.StatusNoContent)
}