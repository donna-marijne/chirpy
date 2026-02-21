package auth_test

import (
	"strings"
	"testing"

	"github.com/donnamarijne/chirpy/internal/auth"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		password string
		want     string
		wantErr  bool
	}{
		{
			name:     "normal case",
			password: "humongous",
			want:     "$argon2id",
			wantErr:  false,
		},
		{
			name:     "empty string",
			password: "",
			want:     "$argon2id",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.HashPassword(tt.password)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("HashPassword() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("HashPassword() succeeded unexpectedly")
			}
			if !strings.HasPrefix(got, tt.want) {
				t.Errorf("HashPassword() = %v, should have prefix %v", got, tt.want)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		password string
		hash     string
		want     bool
		wantErr  bool
	}{
		{
			name:     "matching case",
			password: "humongous",
			hash:     "$argon2id$v=19$m=65536,t=1,p=16$hzZM+uyyLI/jON7LO3EP9A$E7vnbP+T+Aa7wTab7vNbs8uCs6X6ZkkCy8iwFRk+QlQ",
			want:     true,
			wantErr:  false,
		},
		{
			name:     "non-matching case",
			password: "Humongous",
			hash:     "$argon2id$v=19$m=65536,t=1,p=16$hzZM+uyyLI/jON7LO3EP9A$E7vnbP+T+Aa7wTab7vNbs8uCs6X6ZkkCy8iwFRk+QlQ",
			want:     false,
			wantErr:  false,
		},
		{
			name:     "hash unset",
			password: "humongous",
			hash:     "unset",
			want:     false,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.CheckPasswordHash(tt.password, tt.hash)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CheckPasswordHash() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CheckPasswordHash() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
