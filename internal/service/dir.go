package service

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/riazahmedshah/vfs/internal/errs"
	"github.com/riazahmedshah/vfs/internal/model/dir"
	"github.com/riazahmedshah/vfs/internal/model/file"
	"github.com/riazahmedshah/vfs/internal/repository"
	"github.com/riazahmedshah/vfs/internal/server"
)

const (
	msgCreateDirFailed = "failed to create directory"
	msgGetDirFailed    = "failed to get directory"
	msgUpdateDirFailed = "failed to update directory"
)

type DirService struct {
	server   *server.Server
	dirRepo  *repository.DirRepository
	fileRepo *repository.FileRepository
}

func NewDirService(s *server.Server, dirRepo *repository.DirRepository, fileRepo *repository.FileRepository) *DirService {
	return &DirService{
		dirRepo:  dirRepo,
		fileRepo: fileRepo,
	}
}

func (s *DirService) CreateDirectory(ctx context.Context, userID uuid.UUID, parentID uuid.UUID, payload *dir.CreateDirPayload) (*dir.Dir, error) {
	directory, err := s.dirRepo.CreateDirectory(ctx, userID, parentID, payload)
	if err != nil {
		if errors.Is(err, errs.ErrConflictDirName) {
			return nil, err
		}
		return nil, errs.New(http.StatusInternalServerError, msgCreateDirFailed, err)
	}
	return directory, nil
}

func (s *DirService) GetDirectoryContent(ctx context.Context, userID uuid.UUID, dirID uuid.UUID) (*dir.FolderContentResponse, error) {
	var dirData *dir.Dir
	var childDirs []*dir.Dir
	var childFiles []*file.File

	var dirErr, childDirsErr, childFilesErr error
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		dirData, dirErr = s.dirRepo.GetDirectoryById(ctx, userID, dirID)
	}()

	go func() {
		defer wg.Done()
		childDirs, childDirsErr = s.dirRepo.GetChildDirectories(ctx, userID, dirID)
	}()

	go func() {
		defer wg.Done()
		childFiles, childFilesErr = s.fileRepo.GetFilesByDirID(ctx, userID, dirID)
	}()
	wg.Wait()

	if dirErr != nil {
		if errors.Is(dirErr, errs.ErrDirNotFound) {
			return nil, dirErr
		}
		return nil, errs.New(http.StatusInternalServerError, msgGetDirFailed, dirErr)
	}

	if childDirsErr != nil {
		return nil, errs.New(http.StatusInternalServerError, msgGetDirFailed, childDirsErr)
	}
	if childDirs == nil {
		childDirs = []*dir.Dir{}
	}

	if childFilesErr != nil {
		return nil, errs.New(http.StatusInternalServerError, msgGetDirFailed, childFilesErr)
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

func (s *DirService) UpdateDirectory(ctx context.Context, userID uuid.UUID, dirID uuid.UUID, payload *dir.UpdateDirpayload) (*dir.Dir, error) {
	updatedDir, err := s.dirRepo.UpdateDirectory(ctx, userID, dirID, payload)
	if err != nil {
		if errors.Is(err, errs.ErrConflictDirName) {
			return nil, err
		}
		return nil, errs.New(http.StatusInternalServerError, msgUpdateDirFailed, err)
	}
	return updatedDir, nil
}

func (s *DirService) DeleteDirectory(ctx context.Context, userID uuid.UUID, dirID uuid.UUID) error {
	err := s.dirRepo.DeleteDirectory(ctx, userID, dirID)
	if err != nil {
		if errors.Is(err, errs.ErrDirNotFound) {
			return err
		}
		return errs.New(http.StatusInternalServerError, "failed to delete directory", err)
	}
	return nil
}
