package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riazahmedshah/vfs/internal/errs"
	"github.com/riazahmedshah/vfs/internal/model/user"
	"github.com/riazahmedshah/vfs/internal/server"
)

type UserRepository struct {
	server *server.Server
}

func NewUserRepository(s *server.Server) *UserRepository {
	return &UserRepository{
		server: s,
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, tx pgx.Tx, payload *user.CreateUserpayload) (*user.User, error) {
	columns := []string{"username", "email", "image"}
	placeholders := []string{"@username", "@email", "@image"}
	args := pgx.NamedArgs{
		"username": payload.Username,
		"email":    payload.Email,
		"image":    payload.Image,
	}

	if payload.IsGuest {
		columns = append(columns, "is_guest")
		placeholders = append(placeholders, "true")
		args["is_guest"] = true
	}

	if payload.MaxStorageLimit != nil {
		columns = append(columns, "max_storage_limit")
		placeholders = append(placeholders, "@max_storage_limit")
		args["max_storage_limit"] = *payload.MaxStorageLimit
	}

	if payload.MaxFileLimit != nil {
		columns = append(columns, "max_file_limit")
		placeholders = append(placeholders, "@max_file_limit")
		args["max_file_limit"] = *payload.MaxFileLimit
	}

	stmt := fmt.Sprintf(`
		INSERT INTO users (%s)
		VALUES (%s)
		RETURNING *
	`, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	rows, err := tx.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create user query for email=%v: %w", payload.Email, err)
	}

	userItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:users for email=%v: %w", payload.Email, err)
	}

	return &userItem, err
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	stmt := `
		SELECT 
			id, username, email, image, is_guest, max_storage_limit, max_file_limit, created_at, updated_at
		FROM users
		WHERE
			email=@email
	`

	rows, err := u.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"email": email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get user by email query email=%s: %w", email, err)
	}

	userItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:users for email=%s: %w", email, err)
	}

	return &userItem, nil
}

func (u *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*user.ResponseUser, error) {
	stmt := `
		SELECT 
			u.id, u.username, u.email, u.image, u.is_guest, u.max_storage_limit, u.max_file_limit, d.id AS root_dir_id
		FROM users u
		JOIN dirs d ON u.id = d.user_id AND d.parent_id IS NULL
		WHERE
			u.id=@userID
	`
	rows, err := u.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"userID": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get user by id query userID=%s: %w", userID, err)
	}

	userItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.ResponseUser])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to collect row from table:users for userID=%s: %w", userID, err)
	}

	return &userItem, nil
}
