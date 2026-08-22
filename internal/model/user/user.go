package user

import "github.com/riazahmedshah/vfs/internal/model"

type User struct {
	model.Base
	Username        string  `json:"username" db:"username"`
	Email           *string `json:"email" db:"email"`
	Image           *string `json:"image" db:"image"`
	IsGuest         bool    `json:"isGuest" db:"is_guest"`
	MaxStorageLimit int64   `json:"maxStorageLimit" db:"max_storage_limit"`
	MaxFileLimit    int64   `json:"maxFileLimit" db:"max_file_limit"`
}

type UserGetEmail struct {
	Email *string `json:"email" db:"email"`
}
