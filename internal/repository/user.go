package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
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
	stmt := `
		INSERT INTO users (
			username, email, image
		)
		VALUES (
			@username, @email, @image
		)
		RETURNING *
	`
	rows, err := tx.Query(ctx, stmt, pgx.NamedArgs{
		"username": payload.Username,
		"email":    payload.Email,
		"image":    payload.Image,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create user query for email=%s: %w", payload.Email, err)
	}

	userItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:users for email=%s: %w", payload.Email, err)
	}

	return &userItem, err
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	stmt := `
		SELECT 
			id, username, email, image
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
