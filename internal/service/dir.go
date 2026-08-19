package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/riazahmedshah/vfs/internal/errs"
	"github.com/riazahmedshah/vfs/internal/model/dir"
	"github.com/riazahmedshah/vfs/internal/repository"
)

const (
	msgCreateDirFailed = "failed to create directory"
)

type DirService struct {
	dirRepo *repository.DirRepository
}

func NewDirRepository(dirRepo *repository.DirRepository) *DirService {
	return &DirService{
		dirRepo: dirRepo,
	}
}

func (s *DirService) NewDir(ctx context.Context, userID uuid.UUID, payload *dir.CreateDirPayload) (*dir.Dir, error) {
	directory, err := s.dirRepo.CreateDirectory(ctx, userID, payload)
	if err != nil {
		if errors.Is(err, errs.ErrConflictDirName) {
			return nil, err
		}
		return nil, errs.New(http.StatusInternalServerError, msgCreateDirFailed, err)
	}
	return directory, nil
}
