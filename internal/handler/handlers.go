package handler

import (
	"github.com/riazahmedshah/vfs/internal/server"
	"github.com/riazahmedshah/vfs/internal/service"
)

type Handlers struct {
	User *UserHandler
	Dir  *DirHandler
	File *FileHandler
}

func NewHandlers(s *server.Server, service *service.Services) *Handlers {
	return &Handlers{
		User: NewUserHandler(service.User),
		Dir:  NewDirHandler(service.Dir),
		File: NewFileHandler(service.File),
	}
}
