package user

type CreateUserpayload struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Image    *string `json:"image"`
}
