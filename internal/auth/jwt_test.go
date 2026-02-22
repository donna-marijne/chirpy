package auth_test

import (
	"log"
	"testing"
	"time"

	"github.com/donnamarijne/chirpy/internal/auth"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		want        int
		wantErr     bool
	}{
		{
			name:        "normal case",
			userID:      uuid.New(),
			tokenSecret: "humongous",
			expiresIn:   time.Duration(5 * time.Minute),
			want:        1,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("MakeJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MakeJWT() succeeded unexpectedly")
			}
			if len(got) < tt.want {
				t.Errorf("MakeJWT() = %v, want >= %v", got, tt.want)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userA := uuid.New()
	validJWT, err := auth.MakeJWT(userA, "humongous", time.Duration(5*time.Minute))
	if err != nil {
		log.Fatal(err)
	}
	expiredJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiMzc0ZmQ1YTctYTY5MC00OGQwLWFiZjAtYjQ2NWNjNjk0YzA5IiwiZXhwIjoxNzcxNzk5MjQwLCJpYXQiOjE3NzE3OTg5NDB9.9YxtJCBGHfgF6rLAY1_B6Is5bVfX44r3S_1Xvct8zsE"
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tokenString string
		tokenSecret string
		want        uuid.UUID
		wantErr     bool
	}{
		{
			name:        "valid jwt",
			tokenString: validJWT,
			tokenSecret: "humongous",
			want:        userA,
			wantErr:     false,
		},
		{
			name:        "valid jwt with incorrect secret",
			tokenString: validJWT,
			tokenSecret: "Humongous",
			want:        uuid.UUID{},
			wantErr:     true,
		},
		{
			name:        "expired jwt",
			tokenString: expiredJWT,
			tokenSecret: "humongous",
			want:        uuid.UUID{},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.ValidateJWT(tt.tokenString, tt.tokenSecret)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ValidateJWT() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("ValidateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}
