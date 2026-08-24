package dir

import (
	"github.com/google/uuid"
	"github.com/riazahmedshah/vfs/internal/model"
)

type Dir struct {
	model.Base
	Name      string      `json:"name" db:"name"`
	UserID    uuid.UUID   `json:"userId" db:"user_id"`
	ParentID  *uuid.UUID  `json:"parentId" db:"parent_id"`
	Ancestors []uuid.UUID `json:"ancestors" db:"ancestors"`
}

type BreadcrumbItem struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
