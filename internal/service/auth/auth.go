package auth

import "context"

type Identity struct {
	Name    string
	Email   string
	Picture string
}

type Verifier interface {
	Verify(ctx context.Context, token string) (*Identity, error)
}
