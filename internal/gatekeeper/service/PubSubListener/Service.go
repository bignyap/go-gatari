package pubsublistener

import (
	"github.com/bignyap/go-admin/internal/caching"
	gatekeeping "github.com/bignyap/go-admin/internal/gatekeeper/service/GateKeeping"
	"github.com/bignyap/go-utilities/logger/api"
	"github.com/bignyap/go-utilities/pubsub"
)

type PubsubListener struct {
	Logger api.Logger
	Cache  *caching.CacheController
	Match  *gatekeeping.Matcher
	PubSub pubsub.PubSubClient
}

func NewPubSubListener(
	logger api.Logger,
	cache *caching.CacheController,
	matcher *gatekeeping.Matcher,
	pubsubClient pubsub.PubSubClient,
) *PubsubListener {
	return &PubsubListener{
		Logger: logger,
		Cache:  cache,
		Match:  matcher,
		PubSub: pubsubClient,
	}
}
