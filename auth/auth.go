package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetApiKey extract an api key from the header
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authorization info found")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part auth header")
	}
	return vals[1], nil
}
