package main

import (
	"errors"
	"net/http"
)

const (
	userIDHeader = "X-User-ID"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func authenticateUser(req *http.Request) (string, error) {
	var err error
	userID := req.Header.Get(userIDHeader)
	if userID == "" {
		err = ErrUserNotFound
	}

	return userID, err
}
