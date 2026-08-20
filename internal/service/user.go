package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/riazahmedshah/vfs/internal/errs"
	"github.com/riazahmedshah/vfs/internal/lib/utils"
	"github.com/riazahmedshah/vfs/internal/model/user"
	"github.com/riazahmedshah/vfs/internal/repository"
	"github.com/riazahmedshah/vfs/internal/server"
	"github.com/riazahmedshah/vfs/internal/service/auth"
)

const (
	msgCreateUserFailed = "failed to create user"
)

type UserService struct {
	server   *server.Server
	userRepo *repository.UserRepository
	dirRepo  *repository.DirRepository
	verifier auth.Verifier
}

func NewUserService(s *server.Server, userRepo *repository.UserRepository, dirRepo *repository.DirRepository, v auth.Verifier) *UserService {
	return &UserService{
		server:   s,
		userRepo: userRepo,
		dirRepo:  dirRepo,
		verifier: v,
	}
}

func (s *UserService) CreateUser(ctx context.Context, googleToken string) (string, error) {
	identity, err := s.verifier.Verify(ctx, googleToken)
	if err != nil {
		return "", errs.ErrInvalidCredential
	}

	existingUser, err := s.userRepo.GetUserByEmail(ctx, identity.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", errs.New(http.StatusInternalServerError, msgCreateUserFailed, err)
	}
	if err == nil && existingUser.Email == identity.Email {
		token, err := utils.GenerateJWT(existingUser.ID, s.server.Config.Auth.JwtSecret)
		if err != nil {
			return "", errs.New(http.StatusInternalServerError, msgCreateUserFailed, err)
		}
		return token, nil
	}

	tx, err := s.server.DB.Pool.Begin(ctx)
	if err != nil {
		return "", errs.New(http.StatusInternalServerError, msgCreateUserFailed, err)
	}

	defer tx.Rollback(ctx)

	var userPayload user.CreateUserpayload
	userPayload.Username = identity.Name
	userPayload.Email = identity.Email
	userPayload.Image = &identity.Picture
	newUser, err := s.userRepo.CreateUser(ctx, tx, &userPayload)
	if err != nil {
		return "", errs.New(http.StatusInternalServerError, msgCreateUserFailed, err)
	}
	_, err = s.dirRepo.CreateRootDirectory(ctx, tx, newUser.ID)
	if err != nil {
		return "", errs.New(http.StatusInternalServerError, msgCreateUserFailed, err)
	}

	token, err := utils.GenerateJWT(newUser.ID, s.server.Config.Auth.JwtSecret)
	if err != nil {
		return "", errs.New(http.StatusInternalServerError, msgCreateUserFailed, err)
	}

	tx.Commit(ctx)
	return token, nil
}
