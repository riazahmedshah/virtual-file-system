package service

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/riazahmedshah/vfs/internal/errs"
	"github.com/riazahmedshah/vfs/internal/lib/gcs"
	"github.com/riazahmedshah/vfs/internal/model/file"
	"github.com/riazahmedshah/vfs/internal/repository"
	"github.com/riazahmedshah/vfs/internal/server"
)

const (
	msgUploadFileFailed = "failed to upload file"
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

func (s *FileService) UploadFile(ctx context.Context, userID uuid.UUID, reader io.Reader, payload *file.CreateFilePayload) (*file.File, error) {
	randSuffix := uuid.Must(uuid.NewV7())
	objName := fmt.Sprintf("%s/%s-%s", userID, randSuffix, payload.Name)
	if err := s.gcsClient.UploadFile(ctx, objName, reader); err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgUploadFileFailed, err)
	}

	fileItem, err := s.fileRepo.CreateAndUploadFile(ctx, userID, objName, payload)
	if err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgUploadFileFailed, err)
	}

	return fileItem, nil
}
