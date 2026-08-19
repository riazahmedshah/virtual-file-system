package errs

import "fmt"

type AppErr struct {
	Status  int
	Message string
	Err     error
}

func New(status int, message string, internalErr error) *AppErr {
	return &AppErr{
		Status:  status,
		Message: message,
		Err:     internalErr,
	}
}

func (e *AppErr) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppErr) Unwrap() error {
	return e.Err
}
