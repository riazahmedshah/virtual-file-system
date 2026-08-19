package user

import "github.com/riazahmedshah/vfs/internal/model"

type User struct {
	model.Base
	Username string  `json:"username" db:"username"`
	Email    string  `json:"email" db:"email"`
	Image    *string `json:"image" db:"image"`
}
