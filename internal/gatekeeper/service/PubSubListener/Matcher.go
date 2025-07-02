package pubsublistener

import (
	"context"
	"encoding/json"

	"github.com/bignyap/go-admin/internal/common"
	gatekeeping "github.com/bignyap/go-admin/internal/gatekeeper/service/GateKeeping"
	"github.com/bignyap/go-utilities/logger/api"
)

func (s *PubsubListener) UpdateEPMatcher() error {

	if s.PubSub == nil {
		return nil // Pubsub disabled or not configured
	}

	go func() {
		_ = s.PubSub.Subscribe(context.Background(), string(common.EndpointCreated), func(ctx context.Context, payload []byte) error {
			var evt common.EndpointCreatedEvent
			if err := json.Unmarshal(payload, &evt); err != nil {
				s.Logger.Error("failed to unmarshal endpoint event", err)
				return err
			}

			s.Match.Add(gatekeeping.Endpoint{
				Path:   evt.Path,
				Method: evt.Method,
				Code:   evt.Code,
			})

			s.Logger.Info("endpoint matcher updated",
				api.Field{
					Key:   "event",
					Value: evt,
				},
			)
			return nil
		})
	}()

	s.Logger.Info("subscribed to pubsub channel: endpoint:created")
	return nil
}
