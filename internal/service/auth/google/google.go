package google

import (
	"context"
	"fmt"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/riazahmedshah/vfs/internal/config"
	"github.com/riazahmedshah/vfs/internal/service/auth"
)

type Verifier struct {
	cfg *config.Config
}

func NewVerifier(cfg *config.Config) *Verifier {
	return &Verifier{cfg: cfg}
}

func (v *Verifier) Verify(ctx context.Context, token string) (*auth.Identity, error) {
	payload, err := idtoken.Validate(ctx, token, v.cfg.Auth.GoogleClientID)
	if err != nil {
		return nil, fmt.Errorf("invalid google token: %w", err)
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("email claim missing")
	}
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	return &auth.Identity{Email: email, Name: name, Picture: picture}, nil
}
