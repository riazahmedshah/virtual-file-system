package file

type CreateFilePayload struct {
	Name string `form:"name" validate:"required"`
}

type UpdateFilePayload struct {
	Name string `json:"name" validate:"required"`
}
