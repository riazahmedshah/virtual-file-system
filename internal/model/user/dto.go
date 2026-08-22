package user

type CreateUserpayload struct {
	Username        string  `json:"username"`
	Email           *string `json:"email"`
	Image           *string `json:"image"`
	IsGuest         bool    `json:"is_guest"`
	MaxStorageLimit *int64  `json:"max_storage_limit"`
	MaxFileLimit    *int64  `json:"max_file_limit"`
}
