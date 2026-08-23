package errs

import "net/http"

var (
	ErrUserNotFound = New(
		http.StatusNotFound,
		"user not found",
		nil,
	)
)
