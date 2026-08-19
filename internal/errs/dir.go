package errs

import "net/http"

var (
	ErrConflictDirName = New(
		http.StatusConflict,
		"a folder with this name already exists",
		nil,
	)

	ErrInvalidCredential = New(
		http.StatusUnauthorized,
		"invalid or expired credential",
		nil,
	)
)
