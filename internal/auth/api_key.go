package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("Authorization header not present")
	}

	token, ok := strings.CutPrefix(
		strings.TrimSpace(authorization),
		"ApiKey",
	)
	if !ok {
		return "", errors.New("Authorization header missing ApiKey prefix")
	}

	return strings.TrimSpace(token), nil
}
