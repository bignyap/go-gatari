package gatekeeping

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
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

type Endpoint struct {
	Code   string
	Method string
	Path   string
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
	Method           string `json:"method"`
	Path             string `json:"path"`
	OrganizationName string `json:"organization_name"`
}

type RecordUsageInput struct {
	Method           string `json:"method"`
	Path             string `json:"path"`
	OrganizationName string `json:"organization_name"`
}
