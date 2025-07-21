package pubsublistener

import "github.com/bignyap/go-admin/internal/common"

func (s *PubsubListener) ResetGoAdminCache() error {
	if s.PubSub == nil {
		return nil
	}

	s.asyncSubscribe(
		common.OrganizationModified,
		s.cacheInvalidationHandler(common.OrganizationModified, &common.OrganizationModifiedEvent{}),
	)
	s.asyncSubscribe(
		common.SubscriptionModified,
		s.cacheInvalidationHandler(common.SubscriptionModified, &common.SubscriptionModifiedEvent{}),
	)
	s.asyncSubscribe(
		common.PricingModified,
		s.cacheInvalidationHandler(common.PricingModified, &common.PricingModifiedEvent{}),
	)
	s.asyncSubscribe(
		common.PricingModified,
		s.cacheInvalidationHandler(common.OrgPermissionModified, &common.OrgPermissionModifiedEvent{}),
	)

	s.Logger.Info("subscribed to pubsub channels: org/sub/pricing modified")
	return nil
}
