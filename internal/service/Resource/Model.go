package resource

type RegisterEndpointParams struct {
	Name        string  `form:"name" json:"name" validate:"required"`
	Description *string `form:"description" json:"description"` // optional
}

type RegisterEndpointOutputs struct {
	ID int `json:"id"`
	RegisterEndpointParams
}

type CreateResourceTypeParams struct {
	Name        string  `form:"name" json:"name" validate:"required,min=1"`
	Code        string  `form:"code" json:"code" validate:"required,min=1"`
	Description *string `form:"description" json:"description"`
}

type CreateResourceTypeOutput struct {
	ID int `json:"id"`
	CreateResourceTypeParams
}
