package dir

import "github.com/google/uuid"

type CreateDirPayload struct {
	Name     string    `json:"name" validate:"required,min=2,max=255"`
	ParentID uuid.UUID `json:"parentId" validate:"required,uuid"`
}

type UpdateDirpayload struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}
