package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.hash, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secretKey := "secret"
	validToken, errMakeJWT := MakeJWT(userID, secretKey, time.Hour)
	// This should fail because of expiration
	// validToken, errMakeJWT := MakeJWT(userID, "secret", time.Nanosecond)
	if errMakeJWT != nil {
		t.Errorf("couldn't make JWT")
	}

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: secretKey,
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: secretKey,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetAuthToken(t *testing.T) {

	token := "ksC14DD31YGFymjHO/beirK3N/fThRDE/IQfEUDUisz+FfmjvyK1VVkrbsJu1JESAUSfd1FQLOBf1qhIANCCVQ=="
	headerCorrect := make(http.Header)
	headerCorrect.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	headerIncorrect := make(http.Header)
	headerIncorrect.Add("Authorization", "Bearer")
	headerDirty := make(http.Header)
	headerDirty.Add("Authorization", fmt.Sprintf("Bearer                  \n\n%s\n\n", token))
	headerMissing := make(http.Header)

	tests := []struct {
		name    string
		header  http.Header
		token   string
		wantErr bool
	}{
		{
			name:    "Correct Authorization Header",
			header:  headerCorrect,
			token:   token,
			wantErr: false,
		},
		{
			name:    "Authorization Header without token",
			header:  headerIncorrect,
			token:   "",
			wantErr: true,
		},
		{
			name:    "Header with extra spaces and jumps of lines",
			header:  headerDirty,
			token:   token,
			wantErr: false,
		},
		{
			name:    "Missing Header",
			header:  headerMissing,
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GeatBearerToken(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotToken != tt.token {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, token)
			}
		})
	}
}
