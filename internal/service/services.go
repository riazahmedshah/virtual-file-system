package service

import (
	"github.com/riazahmedshah/vfs/internal/repository"
	"github.com/riazahmedshah/vfs/internal/server"
)

type Services struct{}

func NewServices(s *server.Server, r *repository.Repositories) *Services {
	return &Services{}
}
