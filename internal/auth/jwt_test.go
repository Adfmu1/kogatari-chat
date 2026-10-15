package auth

import (
	"errors"
	"testing"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

var testSecretToken = "secret_shhh"

func TestJWTGeneration(t *testing.T) {
	size := 5
	uuIDs := make([]uuid.UUID, size)
	testSecretToken := "secret_shhh"
	var err error
	for i := range size {
		uuIDs[i] = uuid.New()
	}

	tokens := make([]string, size)
	for i := range size {
		tokens[i], err = GenerateJWT(uuIDs[i], testSecretToken, time.Minute)
		if err != nil {
			t.Errorf("GenerateJWT() error %s", err.Error())
		}
	}

	for i := range size {
		uId, err := ValidateJWT(tokens[i], testSecretToken)
		if err != nil {
			t.Errorf("ValidateJWT() error %s", err.Error())
		}
		if uId != uuIDs[i] {
			t.Errorf("ValidateJWT() = %v, want %v", uId, uuIDs[i])
		}
	}
}

func TestExpiredJWT(t *testing.T) {
	testUUID := uuid.New()

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    "kogatari-access",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * -10)),
		Subject:   testUUID.String(),
	})

	jwtToken, err := token.SignedString([]byte(testSecretToken))
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateJWT(jwtToken, testSecretToken)
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Fatal(err)
	}
}

func TestBadIssuerJWT(t *testing.T) {
	testUUID := uuid.New()

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    "bad-issuer",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 10)),
		Subject:   testUUID.String(),
	})

	jwtToken, err := token.SignedString([]byte(testSecretToken))
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateJWT(jwtToken, testSecretToken)
	if !errors.Is(err, jwt.ErrTokenInvalidIssuer) {
		t.Fatal(err)
	}
}

func TestBadSigningMethodJWT(t *testing.T) {
	testUUID := uuid.New()

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, &jwt.RegisteredClaims{
		Issuer:    "bad-issuer",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 10)),
		Subject:   testUUID.String(),
	})

	jwtToken, err := token.SignedString([]byte(testSecretToken))
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateJWT(jwtToken, testSecretToken)
	if errors.Is(jwt.ErrInvalidKeyType, err) {
		t.Fatal(err)
	}
}

func TestNilIssuerJWT(t *testing.T) {
	testUUID := uuid.Nil()

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    "bad-issuer",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 10)),
		Subject:   testUUID.String(),
	})

	jwtToken, err := token.SignedString([]byte(testSecretToken))
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateJWT(jwtToken, testSecretToken)
	if errors.Is(jwt.ErrTokenInvalidIssuer, err) {
		t.Fatal(err)
	}
}
