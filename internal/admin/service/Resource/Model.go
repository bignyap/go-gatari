package resource

type RegisterEndpointParams struct {
	Name           string  `form:"name" json:"name" validate:"required"`
	Description    *string `form:"description" json:"description"`
	HttpMethod     string  `form:"http_method" json:"http_method" validate:"required,oneof=GET POST PUT DELETE PATCH"`
	PathTemplate   string  `form:"path_template" json:"path_template" validate:"required"`
	ResourceTypeID int32   `form:"resource_type_id" json:"resource_type_id" validate:"required"`
}

type RegisterEndpointOutputs struct {
	ID int `json:"id"`
	RegisterEndpointParams
}

type ListEndpointOutputs struct {
	ID               int    `json:"id"`
	ResourceTypeName string `json:"resource_type_name"`
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
