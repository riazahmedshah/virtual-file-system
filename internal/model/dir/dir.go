package dir

import (
	"github.com/google/uuid"
	"github.com/riazahmedshah/vfs/internal/model"
	"github.com/riazahmedshah/vfs/internal/model/file"
)

type Dir struct {
	model.Base
	Name     string     `json:"name" db:"name"`
	UserID   uuid.UUID  `json:"userId" db:"user_id"`
	ParentID *uuid.UUID `json:"parentId" db:"parent_id"`
}

type FolderContentResponse struct {
	CurrentFolder *Dir        `json:"currentFolder"`
	Directories   []Dir       `json:"directories"`
	Files         []file.File `json:"files"`
}
