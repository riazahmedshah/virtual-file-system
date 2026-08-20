package file

import "github.com/google/uuid"

type CreateFilePayload struct {
	Name  string    `form:"name" validate:"required"`
	DirID uuid.UUID `form:"dirId" validate:"required"`
}

type UpdateFilePayload struct {
	Name string `json:"name" validate:"required"`
}
