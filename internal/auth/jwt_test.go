package auth_test

import (
	"testing"
	"time"

	"Kuroashi1995/chirpy/internal/auth"

	"github.com/google/uuid"
)

func TestCreateJWT(t *testing.T) {
	type params struct {
		userID	uuid.UUID
		tokenSecret	string
		expiresIn		time.Duration
	}

	testParams := params{
		userID: uuid.New(),
		tokenSecret: "TestSecret",
		expiresIn: 5 * time.Second,
	}
	_, err := auth.MakeJWT(testParams.userID, testParams.tokenSecret, testParams.expiresIn)
	if err != nil {
		t.Errorf("Token creation failed: %v", err.Error())
	}
}

//VALIDATION TESTS
//ValidateJWT params
type TestParams struct {
	tokenSecret	string
	expiration	time.Duration
}

func TestValidation(t *testing.T) {
	//Initialize cases
	testParams := TestParams {
		tokenSecret: "TestSecret",
		expiration: 5 * time.Second,
	}
	//Validation Ok
	//Token String creation
	userUUID := uuid.New()
	tokenString, err := auth.MakeJWT(userUUID, testParams.tokenSecret, testParams.expiration)
	if err != nil {
		t.Error("Error creating tokenString")
	}
	//Actual testing
	returnedUUID, err := auth.ValidateJWT(tokenString, testParams.tokenSecret)
	if err != nil {
		t.Errorf("Error validating JWT: %v", err.Error())
	} else if returnedUUID != userUUID {
		t.Error("Incorrect returned value")
	}
}

func TestInvalidSignature(t *testing.T) {
	testParams := TestParams {
		tokenSecret: "RealKey",
		expiration: 5 * time.Second,
	}

	userUUID := uuid.New()
	tokenString, err := auth.MakeJWT(userUUID, testParams.tokenSecret, testParams.expiration)
	if err != nil {
		t.Error("Error creating tokenString")
	}
	// Actual testing
	_, err = auth.ValidateJWT(tokenString, "NotRealSecret")
	if err == nil {
		t.Error("Token should not be valid")
	}
}

func TestExpiredValidation(t *testing.T) {
	testparams := TestParams {
		tokenSecret: "NormalKey",
		expiration: 1 * time.Millisecond,
	}
	userUUID := uuid.New()
	tokenString, err := auth.MakeJWT(userUUID, testparams.tokenSecret, testparams.expiration)
	if err != nil {
		t.Error("Error creating tokenString")
	}

	//Actual testing
	time.Sleep(1 * time.Millisecond)
	_, err = auth.ValidateJWT(tokenString, testparams.tokenSecret)
	if err == nil {
		t.Error("This should fail as expired token")
	}
}

