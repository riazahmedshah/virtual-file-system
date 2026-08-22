package file

type CreateFilePayload struct {
	Name string `validate:"min=1,max=255,filename"`
	Ext  string `validate:"required,alphanum,min=1,max=10"`
	Size int64  `validate:"required"`
}

type UpdateFilePayload struct {
	Name string `json:"name" validate:"required"`
}
