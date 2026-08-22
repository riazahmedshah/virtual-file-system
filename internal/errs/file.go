package errs

import "net/http"

var (
	ErrFileNotFound = New(
		http.StatusNotFound,
		"file not found",
		nil,
	)

	ErrMIMETypeMismatch = New(
		http.StatusBadRequest,
		"file MIME type does not match the claimed extension",
		nil,
	)
)
