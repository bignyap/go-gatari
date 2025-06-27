package gatekeeping

import (
	"context"
)

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

type Repository interface {
	GetOrganizationByName(ctx context.Context, name string) (Organization, error)
	GetEndpointByName(ctx context.Context, name string) (ApiEndpoint, error)
	GetActiveSubscription(ctx context.Context, orgID, endpointID int32) (Subscription, error)
	GetPricing(ctx context.Context, subID, endpointID int32) (Pricing, error)
	IncrementUsage(ctx context.Context, subID, endpointID, orgID int32) error
}
