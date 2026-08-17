package repository

import "github.com/riazahmedshah/vfs/internal/server"

type Repositories struct{}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{}
}
