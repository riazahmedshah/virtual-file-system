package service

import (
	"context"
	"fmt"

	"github.com/riazahmedshah/vfs/internal/model/user"
	"github.com/riazahmedshah/vfs/internal/repository"
	"github.com/riazahmedshah/vfs/internal/server"
	"github.com/riazahmedshah/vfs/internal/service/auth"
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

func (s *UserService) CreateUser(ctx context.Context, token string) (*user.User, error) {
	_, err := s.verifier.Verify(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err) // TODO: sentinal err
	}
	return nil, nil
}
