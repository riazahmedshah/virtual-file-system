package dir

import (
	"github.com/riazahmedshah/vfs/internal/model/file"
)

type CreateDirPayload struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

type UpdateDirpayload struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

type FolderContentResponse struct {
	*Dir
	Directories []*Dir       `json:"directories"`
	Files       []*file.File `json:"files"`
}
