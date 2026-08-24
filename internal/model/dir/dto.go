package dir

import (
	"github.com/riazahmedshah/vfs/internal/model/file"
)

type CreateDirPayload struct {
	Name string `json:"name" validate:"required,min=2,max=255,excludesall=/\n\r\t"`
}

type UpdateDirpayload struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

type FolderContentResponse struct {
	*DirResponse
	Directories []*Dir       `json:"directories"`
	Files       []*file.File `json:"files"`
}

type DirResponse struct {
	*Dir
	Breadcrumbs []BreadcrumbItem `json:"breadcrumbs"`
	Size        int64            `json:"size" db:"size"`
}
