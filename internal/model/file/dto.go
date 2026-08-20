package file

import "github.com/google/uuid"

type CreateFilePayload struct {
	Name     string    `json:"name" validate:"required"`
	ParentID uuid.UUID `json:"parentId" validate:"required"`
}
