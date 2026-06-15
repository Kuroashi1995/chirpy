package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	//Configure claims
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	}

	//create token
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//sign the token and return it
	return userToken.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Print(err)
		return uuid.UUID{}, err
	} else if subject, err := token.Claims.GetSubject(); err != nil {
		log.Print(err)
		return uuid.UUID{}, err
	} else if userUUID, err := uuid.Parse(subject); err != nil {
		log.Print(err)
		return uuid.UUID{}, err
	} else {
		return userUUID, nil
	}
}

func GetBearerToken(header http.Header) (string, error) {
	// Get auth header
	authorizationString := header.Get("Authorization")

	// check existance
	if authorizationString == "" {
		log.Print("No authorization header for request")
		return "", errors.New("no Authorization header in request")
	}
	// check correct header format
	splittedAuth := strings.Split(authorizationString, " ")
	if splittedAuth[0] != "BEARER" {
		log.Print("Invalid Authorization header")
		return "", errors.New("invalid authorization header")
	}
	return splittedAuth[1], nil
}
