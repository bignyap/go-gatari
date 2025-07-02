package common

type PubSubChannel string

const (
	EndpointCreated PubSubChannel = "endpoint:created"
)

type EndpointCreatedEvent struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Code   string `json:"code"`
}
