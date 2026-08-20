package errs

import "net/http"

var (
	ErrFileNotFound = New(
		http.StatusNotFound,
		"file not found",
		nil,
	)
)
