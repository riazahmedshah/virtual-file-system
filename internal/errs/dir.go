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

	ErrParentIDRequired = New(
		http.StatusBadRequest,
		"parent id cannot be null",
		nil,
	)

	ErrDirNotFound = New(
		http.StatusNotFound,
		"directory not found or you do not have permission to access it",
		nil,
	)
)
