package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

func HashPassword(password string) (string, error) {

	hashedPwd, errHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errHash != nil {
		return "", fmt.Errorf("couldn't hash password: %w", errHash)
	}
	return string(hashedPwd), nil
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	signingKey := []byte(tokenSecret)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	return newToken.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	claims := jwt.RegisteredClaims{}
	jwtToken, errParse := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if errParse != nil {
		return uuid.Nil, fmt.Errorf("couldn't parse token: %w", errParse)
	}

	id, errGetID := jwtToken.Claims.GetSubject()
	if errGetID != nil {
		return uuid.Nil, fmt.Errorf("couldn't get subject's id: %w", errGetID)
	}

	issuer, errIssuer := jwtToken.Claims.GetIssuer()
	if errIssuer != nil {
		return uuid.Nil, fmt.Errorf("couldn't get issuer: %w", errIssuer)
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("couldn't parse id: %w", err)
	}

	return parsedID, nil
}

func GeatBearerToken(headers http.Header) (string, error) {
	authString := headers.Get("Authorization")
	if authString == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authString, "Bearer"))
	if tokenString == "" {
		return "", errors.New("malformed authorization header")
	}

	return tokenString, nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	rand.Read(key)
	encodedKey := hex.EncodeToString(key)
	return encodedKey, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authString := headers.Get("Authorization")
	if authString == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	apiKey := strings.TrimSpace(strings.TrimPrefix(authString, "ApiKey"))
	if apiKey == "" {
		return "", errors.New("malformed authorization header")
	}

	return apiKey, nil
}
