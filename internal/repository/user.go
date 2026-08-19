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

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:users for email=%s: %w", payload.Email, err)
	}

	return &user, err
}
