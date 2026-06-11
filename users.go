package main

import (
	"Kuroashi1995/chirpy/internal/auth"
	"Kuroashi1995/chirpy/internal/database"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) registerUser(w http.ResponseWriter, r *http.Request) {
	type validUser struct {
		Email *string `json:"email"`
		Password *string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	requestBody := validUser{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if requestBody.Email == nil || requestBody.Password == nil {
		respondWithError(w, http.StatusBadRequest, "Email & Password fields are required")
		return
	}

	hashedPassword, err := auth.HashPassword(*requestBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	createUserParams := database.CreateUserParams{
		Email: *requestBody.Email,
		HashedPassword: hashedPassword,
	}
	createdUser, err := cfg.dbQueries.CreateUser(r.Context(), createUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
	}
	respondWithJson(w, http.StatusCreated, createdUser)
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	//struct for valid request
	type validLogin struct {
		Email *string `json:"email"`
		Password *string `json:"password"`
	}

	//decode request
	decoder := json.NewDecoder(r.Body)
	requestBody := validLogin{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if requestBody.Email == nil || requestBody.Password == nil {
		respondWithError(w, http.StatusBadRequest, "Email & Password fields are required")
		return
	}

	//getting user
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), *requestBody.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	//checking password hash
	passwordMatch, err := auth.CheckPasswordHash(*requestBody.Password, user.HashedPassword )
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	if passwordMatch {
		sanitizedUser := sanitizeUser(user)
		respondWithJson(w, http.StatusOK, sanitizedUser)
		return
	} else {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
}
