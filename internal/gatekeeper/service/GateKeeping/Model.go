package gatekeeping

type Organization struct {
	ID   int32
	Name string
}

type ApiEndpoint struct {
	ID   int32
	Name string
}

type Subscription struct {
	ID              int32
	OrganizationID  int32
	ApiLimit        int32
	ExpiryTimestamp int64
	Active          bool
}

type Pricing struct {
	CostPerCall float64
}

type ValidateRequestInput struct {
	OrganizationName string `json:"organization_name"`
	EndpointName     string `json:"endpoint_name"`
}

type RecordUsageInput struct {
	OrganizationName string `json:"organization_name"`
	EndpointName     string `json:"endpoint_name"`
}
