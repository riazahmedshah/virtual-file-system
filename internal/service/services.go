package service

import (
	"fmt"

	"github.com/riazahmedshah/vfs/internal/lib/gcs"
	"github.com/riazahmedshah/vfs/internal/repository"
	"github.com/riazahmedshah/vfs/internal/server"
	"github.com/riazahmedshah/vfs/internal/service/auth/google"
)

type Services struct {
	User *UserService
	File *FileService
}

func NewServices(s *server.Server, r *repository.Repositories) (*Services, error) {
	gcsClient, err := gcs.NewGCSClient(s.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}
	googleVarifier := google.NewVerifier(s.Config)
	return &Services{
		User: NewUserService(s, r.User, r.Dir, googleVarifier),
		File: NewFileService(s, r.File, gcsClient),
	}, nil
}
