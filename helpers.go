package main

import (
	"Kuroashi1995/chirpy/internal/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	fmt.Printf("Responding with error\n")
	type errorBody struct {
		Error string `json:"error"`
	}
	w.WriteHeader(code)
	response := errorBody{
		Error: message,
	}
	dat, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error marshalling response")
	}
	w.Write(dat)
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	fmt.Printf("Responding with json\n")
	w.WriteHeader(code)
	dat, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
	}
	w.Write(dat)
}

func replaceForbidden(input string) string {
	forbidden := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(input, " ")
	for i, word := range words {
		_, ok := forbidden[strings.ToLower(word)]
		if ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

type SanitizedUser struct {
	ID	uuid.UUID	`json:"id"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	Email	string	`json:"email"`
}

func sanitizeUser(user database.User) SanitizedUser  {
	sanitized := SanitizedUser{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}
	return sanitized
}
