package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) validationHandler(w http.ResponseWriter, r *http.Request) {
	type validationRequest struct {
		Body   *string `json:"body"`
		UserId *string `json:"user_id"`
	}
	type validationError struct {
		Error string `json:"error"`
	}
	type validationSuccess struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	requestBody := validationRequest{}
	log.Printf("DEBUG: decoding\n")
	err := decoder.Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	log.Printf("DEBUG: Checking body field existance\n")
	if requestBody.Body == nil {
		log.Printf("Got into nil body\n")
		respondWithError(w, http.StatusBadRequest, "Body field is required")
		return
	}
	log.Printf("DEBUG: assigning len\n")
	chirp_len := len(*requestBody.Body)
	log.Printf("DEBUG: checking len: %v\n", chirp_len)

	if chirp_len > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
	} else {
		response := validationSuccess{
			CleanedBody: replaceForbidden(*requestBody.Body),
		}
		respondWithJson(w, http.StatusOK, response)
	}

}
