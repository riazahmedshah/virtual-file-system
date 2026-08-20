package repository

import "github.com/riazahmedshah/vfs/internal/server"

type Repositories struct {
	User *UserRepository
	Dir  *DirRepository
	File *FileRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		User: NewUserRepository(s),
		Dir:  NewDirRepository(s),
		File: NewFileRepository(s),
	}
}
