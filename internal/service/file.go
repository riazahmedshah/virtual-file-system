package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/riazahmedshah/vfs/internal/errs"
	"github.com/riazahmedshah/vfs/internal/lib/gcs"
	"github.com/riazahmedshah/vfs/internal/lib/utils"
	"github.com/riazahmedshah/vfs/internal/model/file"
	"github.com/riazahmedshah/vfs/internal/repository"
	"github.com/riazahmedshah/vfs/internal/server"
)

const (
	msgUploadFileFailed   = "failed to upload file"
	msgDeleteFileFailed   = "failed to delete file"
	msgGetSignedURLFailed = "failed to load file"
	msgPermissionDenied   = "user does not have permission to access this file"
)

type FileService struct {
	server    *server.Server
	fileRepo  *repository.FileRepository
	gcsClient *gcs.GCSClient
}

func NewFileService(s *server.Server, fileRepo *repository.FileRepository, client *gcs.GCSClient) *FileService {
	return &FileService{
		fileRepo:  fileRepo,
		gcsClient: client,
	}
}

func (s *FileService) UploadFile(ctx context.Context, userID uuid.UUID, dirID uuid.UUID, src multipart.File, payload *file.CreateFilePayload) (*file.File, error) {
	_, err := utils.DetectAndValidateExt(src, payload.Ext)
	if err != nil {
		return nil, errs.ErrMIMETypeMismatch
	}
	randSuffix := uuid.Must(uuid.NewV7())
	objName := fmt.Sprintf("%s/%s-%s", userID, randSuffix, payload.Name)

	if err := s.gcsClient.UploadFile(ctx, objName, src); err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgUploadFileFailed, err)
	}

	fileItem, err := s.fileRepo.CreateFile(ctx, userID, dirID, objName, payload)
	if err != nil {
		if delErr := s.gcsClient.DeleteFile(ctx, objName); delErr != nil {
			slog.Error("failed to rollback gcs upload after db failure", "error", delErr, "object", objName)
		}
		return nil, errs.New(http.StatusInternalServerError, msgUploadFileFailed, err)
	}

	return fileItem, nil
}

func (s *FileService) DeleteFile(ctx context.Context, userID uuid.UUID, fileID uuid.UUID) error {
	fileItem, err := s.fileRepo.GetFileByID(ctx, userID.String(), fileID)
	if err != nil {
		return errs.New(http.StatusInternalServerError, msgDeleteFileFailed, err)
	}

	if err := s.gcsClient.DeleteFile(ctx, fileItem.GCSKey); err != nil {
		return errs.New(http.StatusInternalServerError, msgDeleteFileFailed, err)
	}

	if err := s.fileRepo.DeleteFile(ctx, userID.String(), fileID.String()); err != nil {
		return errs.New(http.StatusInternalServerError, msgDeleteFileFailed, err)
	}

	return nil
}

func (s *FileService) GetFileDetailsByID(ctx context.Context, userID uuid.UUID, fileID uuid.UUID) (*file.File, error) {
	fileItem, err := s.fileRepo.GetFileByID(ctx, userID.String(), fileID)
	if err != nil {
		return nil, errs.New(http.StatusInternalServerError, "failed to get file details", err)
	}

	return fileItem, nil
}

func (s *FileService) GenerateSignedURL(ctx context.Context, userID uuid.UUID, fileID uuid.UUID, disposition gcs.Disposition, expiry time.Duration) (string, error) {
	fileItem, err := s.fileRepo.GetFileByID(ctx, userID.String(), fileID)
	if err != nil {
		if errors.Is(err, errs.ErrFileNotFound) {
			return "", err
		}
		return "", errs.New(http.StatusInternalServerError, msgGetSignedURLFailed, err)
	}

	if fileItem.UserID != userID {
		return "", errs.New(http.StatusForbidden, msgPermissionDenied, nil)
	}
	signedURL, err := s.gcsClient.GenerateSignedURL(fileItem.GCSKey, disposition, expiry)
	if err != nil {

		return "", errs.New(http.StatusInternalServerError, msgGetSignedURLFailed, err)
	}

	return signedURL, nil
}

func (s *FileService) GetUserTotalStorageUsed(ctx context.Context, userID uuid.UUID) (int64, error) {
	totalUsed, err := s.fileRepo.GetUserTotalStorageUsed(ctx, userID)
	if err != nil {
		return 0, errs.New(http.StatusInternalServerError, "failed to get user total storage used", err)
	}
	return totalUsed, nil
}
