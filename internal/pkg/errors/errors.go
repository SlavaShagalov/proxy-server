package errors

import (
	"errors"
)

var (
	// Common repository
	ErrDb = errors.New("db error")

	// Requests
	ErrRequestNotFound = errors.New("request not found")
)
