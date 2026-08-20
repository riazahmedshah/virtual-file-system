package file

import (
	"github.com/google/uuid"
	"github.com/riazahmedshah/vfs/internal/model"
)

type File struct {
	model.Base
	Name   string    `json:"name" db:"name"`
	DirID  uuid.UUID `json:"dirId" db:"dir_id"`
	UserID uuid.UUID `json:"userId" db:"user_id"`
	GCSKey string    `json:"gcsKey" db:"gcs_key"`
}
