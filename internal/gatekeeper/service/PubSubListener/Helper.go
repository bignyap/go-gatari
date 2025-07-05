package pubsublistener

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bignyap/go-admin/internal/common"
	"github.com/bignyap/go-utilities/logger/api"
)

func (s *PubsubListener) asyncSubscribe(key common.PubSubChannel, handler func(context.Context, []byte) error) {
	go func() {
		_ = s.PubSub.Subscribe(context.Background(), string(key), handler)
	}()
}

func (s *PubsubListener) logUnmarshalError(key common.PubSubChannel, err error) error {
	s.Logger.Error(fmt.Sprintf("failed to unmarshal %s", key), err)
	return err
}

func (s *PubsubListener) cacheInvalidationHandler(key common.PubSubChannel, evtPtr any) func(context.Context, []byte) error {
	return func(ctx context.Context, payload []byte) error {
		if err := json.Unmarshal(payload, evtPtr); err != nil {
			return s.logUnmarshalError(key, err)
		}

		id := extractID(evtPtr)
		s.Cache.DeleteRedisValue(ctx, string(common.OrganizationPrefix), fmt.Sprintf("%s:%s", id, "*"))
		s.Logger.Info("cache removed", api.Field{Key: "event", Value: evtPtr})
		return nil
	}
}

func extractID(evt any) string {
	switch e := evt.(type) {
	case *common.OrganizationModifiedEvent:
		return strconv.Itoa(int(e.ID))
	case *common.SubscriptionModifiedEvent:
		return strconv.Itoa(int(e.ID))
	case *common.PricingModifiedEvent:
		return strconv.Itoa(int(e.ID))
	default:
		return "unknown"
	}
}
