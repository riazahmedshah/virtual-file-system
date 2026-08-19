package service

import (
	"context"
	"net/http"

	"github.com/riazahmedshah/vfs/internal/errs"
	"github.com/riazahmedshah/vfs/internal/model/dir"
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

func (s *UserService) CreateUser(ctx context.Context, token string) (any, error) {
	identity, err := s.verifier.Verify(ctx, token)
	if err != nil {
		return nil, errs.ErrInvalidCredential
	}
	return identity, nil
	tx, err := s.server.DB.Pool.Begin(ctx)
	if err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgCreateUserFailed, err)
	}
	defer tx.Rollback(ctx)

	var userPayload user.CreateUserpayload
	userPayload.Username = identity.Name
	userPayload.Email = identity.Email
	userPayload.Image = &identity.Picture
	user, err := s.userRepo.CreateUser(ctx, tx, &userPayload)
	if err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgCreateUserFailed, err)
	}

	var dirPayload dir.CreateDirPayload
	dirPayload.Name = "root"
	dirPayload.ParentID = nil
	_, err = s.dirRepo.CreateDirectory(ctx, tx, user.ID.String(), &dirPayload)
	if err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgCreateUserFailed, err)
	}

	tx.Commit(ctx)
	return user, nil
}
