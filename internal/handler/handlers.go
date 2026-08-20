package handler

import (
	"github.com/riazahmedshah/vfs/internal/server"
	"github.com/riazahmedshah/vfs/internal/service"
)

type Handlers struct {
	User *UserHandler
	File *FileHandler
}

func NewHandlers(s *server.Server, service *service.Services) *Handlers {
	return &Handlers{
		User: NewUserHandler(service.User),
		File: NewFileHandler(service.File),
	}
}
