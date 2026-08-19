package repository

import "github.com/riazahmedshah/vfs/internal/server"

type Repositories struct {
	User *UserRepository
	Dir  *DirRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		User: NewUserRepository(s),
		Dir:  NewDirRepository(s),
	}
}
