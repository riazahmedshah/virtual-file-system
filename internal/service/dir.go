package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/riazahmedshah/vfs/internal/errs"
	"github.com/riazahmedshah/vfs/internal/model/dir"
	"github.com/riazahmedshah/vfs/internal/model/file"
	"github.com/riazahmedshah/vfs/internal/repository"
)

const (
	msgCreateDirFailed = "failed to create directory"
	msgGetDirFailed    = "failed to get directory"
)

type DirService struct {
	dirRepo  *repository.DirRepository
	fileRepo *repository.FileRepository
}

func NewDirRepository(dirRepo *repository.DirRepository, fileRepo *repository.FileRepository) *DirService {
	return &DirService{
		dirRepo:  dirRepo,
		fileRepo: fileRepo,
	}
}

func (s *DirService) CreateDirectory(ctx context.Context, userID uuid.UUID, payload *dir.CreateDirPayload) (*dir.Dir, error) {
	directory, err := s.dirRepo.CreateDirectory(ctx, userID, payload)
	if err != nil {
		if errors.Is(err, errs.ErrConflictDirName) {
			return nil, err
		}
		return nil, errs.New(http.StatusInternalServerError, msgCreateDirFailed, err)
	}
	return directory, nil
}

func (s *DirService) GetDirectoryByID(ctx context.Context, userID uuid.UUID, dirID uuid.UUID) (*dir.FolderContentResponse, error) {
	dirData, err := s.dirRepo.GetDirectoryById(ctx, userID, dirID)
	if err != nil {
		if errors.Is(err, errs.ErrDirNotFound) {
			return nil, err
		}
		return nil, errs.New(http.StatusInternalServerError, msgGetDirFailed, err)
	}

	childDirs, err := s.dirRepo.GetChildDirectories(ctx, userID, dirID)
	if err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgGetDirFailed, err)
	}

	if childDirs == nil {
		childDirs = []*dir.Dir{}
	}

	childFiles, err := s.fileRepo.GetFilesByDirID(ctx, userID, dirID)
	if err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgGetDirFailed, err)
	}

	if childFiles == nil {
		childFiles = []*file.File{}
	}

	response := &dir.FolderContentResponse{
		Dir:         dirData,
		Directories: childDirs,
		Files:       childFiles,
	}

	return response, nil
}
