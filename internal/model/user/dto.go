package user

import (
	"github.com/google/uuid"
)

type CreateUserpayload struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Username        string    `json:"username" db:"username"`
	Email           *string   `json:"email" db:"email"`
	Image           *string   `json:"image" db:"image"`
	IsGuest         bool      `json:"isGuest" db:"is_guest"`
	MaxStorageLimit *int64    `json:"maxStorageLimit" db:"max_storage_limit"`
	MaxFileLimit    *int64    `json:"maxFileLimit" db:"max_file_limit"`
}

type ResponseUser struct {
	CreateUserpayload
	RootDirID        uuid.UUID `json:"rootDirId" db:"root_dir_id"`
	TotalStorageUsed int64     `json:"totalStorageUsed" db:"total_storage_used"`
}
