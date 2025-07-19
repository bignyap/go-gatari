package gatekeeping

import (
	"net/http"
	"sync"

	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/julienschmidt/httprouter"
)

type Organization struct {
	ID    int32
	Name  string
	Realm string
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

type Endpoint struct {
	Code   string `json:"code" form:"code"`
	Method string `json:"method" form:"method"`
	Path   string `json:"path" form:"path"`
}

type Matcher struct {
	router    *httprouter.Router
	lock      sync.RWMutex
	endpoints map[string]Endpoint
}

type capture struct {
	header http.Header
	code   string
	found  bool
}

type ValidateRequestInput struct {
	Method           string `json:"method" form:"method"`
	Path             string `json:"path" form:"path"`
	OrganizationName string `json:"organization_name" form:"organization_name"`
}

type ValidationRequestOutput struct {
	Organization sqlcgen.GetOrganizationByNameRow `json:"organization"`
	Endpoint     sqlcgen.GetApiEndpointByNameRow  `json:"endpoint"`
	Subscription sqlcgen.GetActiveSubscriptionRow `json:"subscription"`
	Remaining    int32                            `json:"remaining"` // nil if unlimited
}

type RecordUsageInput struct {
	Method           string `json:"method" form:"method"`
	Path             string `json:"path" form:"path"`
	OrganizationName string `json:"organization_name" form:"organization_name"`
}

type GetOrgSubDetailsOutput struct {
	ValidationRequestOutput
	EndpointCode string `json:"endpoint_code"`
}
