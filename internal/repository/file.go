package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riazahmedshah/vfs/internal/errs"
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

func (r *FileRepository) CreateFile(ctx context.Context, userID uuid.UUID, dirID uuid.UUID, gcsKey string, payload *file.CreateFilePayload) (*file.File, error) {
	stmt := `
		INSERT INTO files (
			name, dir_id, user_id, gcs_key, size, ext	
		)
		VALUES (
			@name, @dir_id, @user_id, @gcs_key, @size, @ext
		)
		RETURNING *
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"name":    payload.Name,
		"dir_id":  dirID,
		"user_id": userID,
		"gcs_key": gcsKey,
		"size":    payload.Size,
		"ext":     payload.Ext,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create file query for user_id=%s parent_id=%s: %w", userID, dirID, err)
	}

	fileItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[file.File])
	if err != nil {
		return nil, fmt.Errorf("failed to collect file row from files for user_id=%s parent_id=%s: %w", userID, dirID, err)
	}

	return &fileItem, nil
}

func (r *FileRepository) GetFileByID(ctx context.Context, userID string, fileID string) (*file.File, error) {
	stmt := `
		SELECT 
			id, name, parent_id, user_id, gcs_key, created_at, updated_at, size, ext
		FROM files
		WHERE 
			id = @file_id AND user_id = @user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"file_id": fileID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get file by id query for user_id=%s file_id=%s: %w", userID, fileID, err)
	}

	fileItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[file.File])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to collect file row from files for user_id=%s file_id=%s: %w", userID, fileID, err)
	}

	return &fileItem, nil
}

func (r *FileRepository) UpdateFile(ctx context.Context, userID string, fileID string, payload *file.UpdateFilePayload) (*file.File, error) {
	stmt := `
		UPDATE files
		SET name = @name, updated_at = NOW()
		WHERE id = @file_id AND user_id = @user_id
		RETURNING *
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"name":    payload.Name,
		"file_id": fileID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute update file query for user_id=%s file_id=%s: %w", userID, fileID, err)
	}

	fileItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[file.File])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to collect file row from files for user_id=%s file_id=%s: %w", userID, fileID, err)
	}

	return &fileItem, nil
}

func (r *FileRepository) DeleteFile(ctx context.Context, userID string, fileID string) error {
	stmt := `
		DELETE FROM files
		WHERE id = @file_id AND user_id = @user_id
	`

	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"file_id": fileID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to execute delete file query for user_id=%s file_id=%s: %w", userID, fileID, err)
	}

	if result.RowsAffected() == 0 {
		return errs.ErrFileNotFound
	}

	return nil
}

func (r *FileRepository) GetFilesByDirID(ctx context.Context, userID uuid.UUID, dirID uuid.UUID) ([]*file.File, error) {
	stmt := `
		SELECT
			id, name, dir_id, user_id, gcs_key, created_at, updated_at
		FROM files
		WHERE dir_id = @dir_id AND user_id = @user_id
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"dir_id":  dirID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get files by dir id query for user_id=%s dir_id=%s: %w", userID, dirID, err)
	}

	fileItems, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[file.File])
	if err != nil {
		return nil, fmt.Errorf("failed to collect file rows from files for user_id=%s dir_id=%s: %w", userID, dirID, err)
	}

	return fileItems, nil
}

func (r *FileRepository) GetUserTotalStorageUsed(ctx context.Context, userID uuid.UUID) (int64, error) {
	stmt := `
		SELECT COALESCE(SUM(size), 0) AS total_used
		FROM files
		WHERE user_id = @user_id
	`
	var totalUsed int64
	err := r.server.DB.Pool.QueryRow(ctx, stmt, pgx.NamedArgs{
		"user_id": userID,
	}).Scan(&totalUsed)
	if err != nil {
		return 0, fmt.Errorf("failed to get total storage used for user_id=%s: %w", userID, err)
	}

	return totalUsed, nil
}
