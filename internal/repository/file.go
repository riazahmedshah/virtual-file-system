package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/riazahmedshah/vfs/internal/model/file"
	"github.com/riazahmedshah/vfs/internal/server"
)

type FileRepository struct {
	server *server.Server
}

func NewFileRepository(s *server.Server) *FileRepository {
	return &FileRepository{
		server: s,
	}
}

func (r *FileRepository) CreateAndUploadFile(ctx context.Context, userID string, gcsKey string, payload *file.CreateFilePayload) (*file.File, error) {
	stmt := `
		INSERT INTO files (
			name, parent_id, user_id, gcs_key	
		)
		VALUES (
			@name, @parent_id, @user_id, @gcs_key
		)
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"name":      payload.Name,
		"parent_id": payload.ParentID,
		"user_id":   userID,
		"gcs_key":   gcsKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create file query for user_id=%s parent_id=%s: %w", userID, payload.ParentID, err)
	}

	fileItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[file.File])
	if err != nil {
		return nil, fmt.Errorf("failed to collect file row from files for user_id=%s parent_id=%s: %w", userID, payload.ParentID, err)
	}

	return &fileItem, nil
}
