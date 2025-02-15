package auth

import (
	"encoding/base64"
	"errors"
	"os"
)

func GetAuthHeaderValue() (string, error) {
	// Get the environment variables
	email := os.Getenv("TOGGL_EMAIL")
	password := os.Getenv("TOGGL_PASSWORD")
	if email == "" {
		return "", errors.New("environment variable TOGGL_EMAIL is not set")
	}
	if password == "" {
		return "", errors.New("environment variable TOGGL_PASSWORD is not set")
	}

	// Concatenate and encode to base64
	auth := email + ":" + password
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))
	headerValue := "Basic " + authEncoded

	return headerValue, nil
}
