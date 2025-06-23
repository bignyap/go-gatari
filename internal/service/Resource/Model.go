package resource

type RegisterEndpointParams struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"` // optional
}

type RegisterEndpointOutputs struct {
	ID int `json:"id"`
	RegisterEndpointParams
}

type CreateResourceTypeParams struct {
	Name        string  `json:"name" validate:"required,min=1"`
	Code        string  `json:"code" validate:"required,min=1"`
	Description *string `json:"description"`
}

type CreateResourceTypeOutput struct {
	ID int `json:"id"`
	CreateResourceTypeParams
}
