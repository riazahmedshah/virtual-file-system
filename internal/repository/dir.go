package repository

import (
	"context"
	"errors"
	"fmt"

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

func (d *DirRepository) CreateDirectory(ctx context.Context, tx pgx.Tx, userID string, paylaod *dir.CreateDirPayload) (*dir.Dir, error) {
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

	dir, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dir.Dir])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:dirs for user_id=%s parent_id=%s %w", userID, paylaod.ParentID, err)
	}

	return &dir, nil
}
