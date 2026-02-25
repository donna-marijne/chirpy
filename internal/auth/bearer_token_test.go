package auth_test

import (
	"net/http"
	"testing"

	"github.com/donnamarijne/chirpy/internal/auth"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		headers http.Header
		want    string
		wantErr bool
	}{
		{
			name: "normal case",
			headers: http.Header{
				"Authorization": []string{"Bearer f00b8r"},
				"Host":          []string{"api.chirpy.com"},
			},
			want:    "f00b8r",
			wantErr: false,
		},
		{
			name: "missing header",
			headers: http.Header{
				"Host": []string{"api.chirpy.com"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "missing prefix",
			headers: http.Header{
				"Authorization": []string{"f00b8r"},
				"Host":          []string{"api.chirpy.com"},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "whitespace prefix",
			headers: http.Header{
				"Authorization": []string{"   Bearer f00b8r"},
				"Host":          []string{"api.chirpy.com"},
			},
			want:    "f00b8r",
			wantErr: false,
		},
		{
			name: "whitespace suffix",
			headers: http.Header{
				"Authorization": []string{"Bearer f00b8r    "},
				"Host":          []string{"api.chirpy.com"},
			},
			want:    "f00b8r",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.GetBearerToken(tt.headers)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBearerToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBearerToken() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
