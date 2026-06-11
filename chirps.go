package main

import (
	"Kuroashi1995/chirpy/internal/database"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) validateAndSaveChirp(w http.ResponseWriter, r *http.Request) {
	type validationRequest struct {
		Body   *string `json:"body"`
		UserID *string `json:"user_id"`
	}
	type validationError struct {
		Error string `json:"error"`
	}
	decoder := json.NewDecoder(r.Body)
	requestBody := validationRequest{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if requestBody.Body == nil || requestBody.UserID == nil {
		respondWithError(w, http.StatusBadRequest, "All body fields are required")
		return
	}
	chirpLen := len(*requestBody.Body)

	if chirpLen > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
	}

	validatedUUID, err := uuid.Parse(*requestBody.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid uuid")
	}

	createChirpParams := database.CreateChirpParams{
		Body:   *requestBody.Body,
		UserID: validatedUUID,
	}

	createdChirp, err := cfg.dbQueries.CreateChirp(r.Context(), createChirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not save chirp")
	}
	respondWithJson(w, http.StatusCreated, createdChirp)
}

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps")
	}
	respondWithJson(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "invalidRequest")
	}
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Cannot parse ID")
		return
	}
	chirp, err := cfg.dbQueries.GetChirpById(r.Context(), parsedUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, chirp)
}
