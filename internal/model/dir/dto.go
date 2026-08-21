package dir

import (
	"github.com/google/uuid"
	"github.com/riazahmedshah/vfs/internal/model/file"
)

type CreateDirPayload struct {
	Name     string    `json:"name" validate:"required,min=2,max=255"`
	ParentID uuid.UUID `json:"parentId" validate:"required,uuid"`
}

type UpdateDirpayload struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

type FolderContentResponse struct {
	*Dir
	Directories []*Dir       `json:"directories"`
	Files       []*file.File `json:"files"`
}
