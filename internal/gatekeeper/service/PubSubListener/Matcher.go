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

	// Listen for endpoint creation
	go func() {
		_ = s.PubSub.Subscribe(context.Background(), string(common.EndpointCreated), func(ctx context.Context, payload []byte) error {
			var evt common.EndpointCreatedEvent
			if err := json.Unmarshal(payload, &evt); err != nil {
				s.Logger.Error("failed to unmarshal endpoint:created", err)
				return err
			}

			s.Match.Add(gatekeeping.Endpoint{
				Path:   evt.Path,
				Method: evt.Method,
				Code:   evt.Code,
			})

			s.Logger.Info("endpoint added to matcher",
				api.Field{Key: "event", Value: evt},
			)
			return nil
		})
	}()

	// Listen for endpoint deletion
	go func() {
		_ = s.PubSub.Subscribe(context.Background(), string(common.EndpointDeleted), func(ctx context.Context, payload []byte) error {
			var evt common.EndpointDeletedEvent
			if err := json.Unmarshal(payload, &evt); err != nil {
				s.Logger.Error("failed to unmarshal endpoint:deleted", err)
				return err
			}

			s.Match.Drop(evt.Code)

			s.Logger.Info("endpoint removed from matcher",
				api.Field{Key: "event", Value: evt},
			)
			return nil
		})
	}()

	s.Logger.Info("subscribed to pubsub channels: endpoint:created, endpoint:deleted")
	return nil
}
