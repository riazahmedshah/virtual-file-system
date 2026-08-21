package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/riazahmedshah/vfs/internal/errs"
	"github.com/riazahmedshah/vfs/internal/model/dir"
	"github.com/riazahmedshah/vfs/internal/server"
)

type DirRepository struct {
	server *server.Server
}

func NewDirRepository(s *server.Server) *DirRepository {
	return &DirRepository{
		server: s,
	}
}

func (r *DirRepository) CreateRootDirectory(ctx context.Context, tx pgx.Tx, userID uuid.UUID) (*dir.Dir, error) {
	stmt := `
		INSERT INTO dirs (
			name, user_id, parent_id
		)
		VALUES (
			@name, @user_id, @parent_id
		)
		RETURNING *
	`

	rows, err := tx.Query(ctx, stmt, pgx.NamedArgs{
		"name":      "root",
		"user_id":   userID,
		"parent_id": nil,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique key violation
			return nil, errs.ErrConflictDirName
		}
		return nil, fmt.Errorf("failed to execute create dir query for user_id=%s : %w", userID, err)
	}

	dirItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dir.Dir])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:dirs for user_id=%s: %w", userID, err)
	}

	return &dirItem, nil
}

func (r *DirRepository) CreateDirectory(ctx context.Context, userID uuid.UUID, paylaod *dir.CreateDirPayload) (*dir.Dir, error) {
	stmt := `
		INSERT INTO dirs (
			name, user_id, parent_id
		)
		VALUES (
			@name, @user_id, @parent_id
		)
		RETURNING *
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"name":      paylaod.Name,
		"user_id":   userID,
		"parent_id": paylaod.ParentID,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique key violation
			return nil, errs.ErrConflictDirName
		}
		return nil, fmt.Errorf("failed to execute create dir query for user_id=%s parent_id=%s : %w", userID, paylaod.ParentID, err)
	}

	dirItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dir.Dir])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:dirs for user_id=%s parent_id=%s %w", userID, paylaod.ParentID, err)
	}

	return &dirItem, nil
}

func (r *DirRepository) GetDirectoryById(ctx context.Context, userID, dirID uuid.UUID) (*dir.Dir, error) {
	stmt := `
		SELECT
			id, name, parent_id, user_id, created_at, updated_at
		FROM dirs
		WHERE id = @dir_id AND user_id = @user_id
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"dir_id":  dirID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get query for id=%s user_id=%s: %w", dirID, userID, err)
	}

	dirItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dir.Dir])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrDirNotFound
		}
		return nil, fmt.Errorf("failed to collect row from table:dirs for id=%s user_id=%s: %w", dirID, userID, err)
	}

	return &dirItem, nil
}

func (r *DirRepository) GetChildDirectories(ctx context.Context, userID, dirID uuid.UUID) ([]*dir.Dir, error) {
	stmt := `
		SELECT
			id, name, parent_id, user_id, created_at, updated_at
		FROM dirs
		WHERE parent_id = @parent_id AND user_id = @user_id
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"parent_id": dirID,
		"user_id":   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get query for parent_id=%s user_id=%s: %w", dirID, userID, err)
	}

	dirItems, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[dir.Dir])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows from table:dirs for parent_id=%s user_id=%s: %w", dirID, userID, err)
	}

	return dirItems, nil
}

func (r *DirRepository) UpdateDirectory(ctx context.Context, userID, dirID uuid.UUID, payload *dir.UpdateDirpayload) (*dir.Dir, error) {
	stmt := `
		UPDATE dirs
		SET
			name=@name
			updated_at = NOW()
		WHERE
			id=@dir_id AND user_id=@user_id
		RETURNING *
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"dir_id":  dirID,
		"user_id": userID,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errs.ErrConflictDirName
		}
		return nil, fmt.Errorf("failed to execute update query for id=%s user_id=%s: %w", dirID, userID, err)
	}

	updatedDir, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dir.Dir])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrDirNotFound
		}
		return nil, fmt.Errorf("failed to collect updated row from table:dirs id=%s user_id=%s: %w", dirID, userID, err)
	}

	return &updatedDir, nil
}

func (r *DirRepository) DeleteDirectory(ctx context.Context, userID, dirID uuid.UUID, payload *dir.UpdateDirpayload) error {
	stmt := `
		DELETE FROM dirs
		WHERE
			id=@dir_id AND user_id=@user_id
	`
	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"dir_id":  dirID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to execute update query for id=%s user_id=%s: %w", dirID, userID, err)
	}

	if result.RowsAffected() == 0 {
		return errs.ErrDirNotFound
	}

	return nil
}
