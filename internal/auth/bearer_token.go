package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("Authorization header not present")
	}

	token, ok := strings.CutPrefix(
		strings.TrimSpace(authorization),
		"Bearer",
	)
	if !ok {
		return "", errors.New("Authorization header missing Bearer prefix")
	}

	return strings.TrimSpace(token), nil
}
