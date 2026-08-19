package service

import (
	"github.com/riazahmedshah/vfs/internal/repository"
	"github.com/riazahmedshah/vfs/internal/server"
	"github.com/riazahmedshah/vfs/internal/service/auth/google"
)

type Services struct {
	User *UserService
}

func NewServices(s *server.Server, r *repository.Repositories) *Services {
	googleVarifier := google.NewVerifier(s.Config)
	return &Services{
		User: NewUserService(s, r.User, r.Dir, googleVarifier),
	}
}
