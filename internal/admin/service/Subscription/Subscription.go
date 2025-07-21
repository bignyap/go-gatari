package subscription

import (
	"context"
	"strings"
	"time"

	"github.com/bignyap/go-admin/internal/common"
	"github.com/bignyap/go-admin/internal/database/dbutils"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/converter"
	"github.com/bignyap/go-utilities/server"
	"github.com/jackc/pgx/v5"
	"github.com/jinzhu/copier"
)

type BulkSubscriptionInserter struct {
	Subscriptions       []sqlcgen.CreateSubscriptionsParams
	SubscriptionService *SubscriptionService
}

func (input BulkSubscriptionInserter) InsertRows(ctx context.Context, tx pgx.Tx) (int64, error) {
	return input.SubscriptionService.DB.CreateSubscriptions(ctx, input.Subscriptions)
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, input *CreateSubscriptionParams) (CreateSubscriptionOutput, error) {

	params := sqlcgen.CreateSubscriptionParams{
		SubscriptionName:               input.Name,
		SubscriptionType:               input.Type,
		SubscriptionCreatedDate:        int32(input.CreatedAt.Unix()),
		SubscriptionUpdatedDate:        int32(input.UpdatedAt.Unix()),
		SubscriptionStartDate:          int32(input.StartDate.Unix()),
		SubscriptionApiLimit:           converter.ToPgInt4(input.APILimit),
		SubscriptionExpiryDate:         converter.ToPgInt4FromTimeOrDate(input.ExpiryDate),
		SubscriptionDescription:        converter.ToPgText(input.Description),
		SubscriptionStatus:             converter.ToPgBool(input.Status),
		OrganizationID:                 int32(input.OrganizationID),
		SubscriptionTierID:             int32(input.SubscriptionTierID),
		SubscriptionBillingInterval:    converter.ToPgText(input.BillingInterval),
		SubscriptionBillingModel:       converter.ToPgText(input.BillingModel),
		SubscriptionQuotaResetInterval: converter.ToPgText(input.QuotaResetInterval),
	}

	insertedID, err := s.DB.CreateSubscription(ctx, params)
	if err != nil {
		return CreateSubscriptionOutput{}, server.NewError(
			server.ErrorInternal,
			"couldn't create subscription",
			err,
		)
	}

	output := CreateSubscriptionOutput{
		ID:                       int(insertedID),
		CreateSubscriptionParams: *input,
	}

	return output, nil
}

func (s *SubscriptionService) CreateSubscriptionInBatch(ctx context.Context, inputs []CreateSubscriptionParams) (int, error) {

	var params []sqlcgen.CreateSubscriptionsParams
	currentTime := int32(time.Now().Unix())

	for _, input := range inputs {
		if err := s.Validator.Struct(input); err != nil {
			return 0, server.NewError(
				server.ErrorInternal,
				"validation failed",
				err,
			)
		}

		params = append(params, sqlcgen.CreateSubscriptionsParams{
			SubscriptionName:               input.Name,
			SubscriptionType:               input.Type,
			SubscriptionCreatedDate:        currentTime,
			SubscriptionUpdatedDate:        currentTime,
			SubscriptionStartDate:          currentTime,
			SubscriptionApiLimit:           converter.ToPgInt4(input.APILimit),
			SubscriptionExpiryDate:         converter.ToPgInt4FromTimeOrDate(input.ExpiryDate),
			SubscriptionDescription:        converter.ToPgText(input.Description),
			SubscriptionStatus:             converter.ToPgBool(input.Status),
			OrganizationID:                 int32(input.OrganizationID),
			SubscriptionTierID:             int32(input.SubscriptionTierID),
			SubscriptionBillingInterval:    converter.ToPgText(input.BillingInterval),
			SubscriptionBillingModel:       converter.ToPgText(input.BillingModel),
			SubscriptionQuotaResetInterval: converter.ToPgText(input.QuotaResetInterval),
		})
	}

	inserter := BulkSubscriptionInserter{
		Subscriptions:       params,
		SubscriptionService: s,
	}

	affectedRows, err := dbutils.InsertWithTransaction(ctx, s.Conn, inserter)
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal,
			"insert failed",
			err,
		)
	}

	return int(affectedRows), nil
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, idType string, Id int) error {

	var pusbSubId int32

	switch strings.ToLower(idType) {
	case "organization":
		err := s.DB.DeleteSubscriptionByOrgId(ctx, int32(Id))
		if err != nil {
			return server.NewError(
				server.ErrorInternal,
				"couldn't delete the subscription",
				err,
			)
		}
		pusbSubId = int32(Id)
	case "subscription":
		orgId, err := s.DB.DeleteSubscriptionById(ctx, int32(Id))
		if err != nil {
			return server.NewError(
				server.ErrorInternal,
				"couldn't delete the subscription",
				err,
			)
		}
		pusbSubId = orgId
	}

	err := s.PubSubClient.Publish(ctx, string(common.SubscriptionModified), common.SubscriptionModifiedEvent{
		ID: pusbSubId,
	})
	if err != nil {
		return server.NewError(
			server.ErrorInternal,
			"couldn't push to the queue",
			err,
		)
	}

	return nil
}

func (s *SubscriptionService) GetSubscription(ctx context.Context, id int) (ListSubscriptionOutput, error) {

	subscription, err := s.DB.GetSubscriptionById(ctx, int32(id))
	if err != nil {
		return ListSubscriptionOutput{}, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the subscription",
			err,
		)
	}

	row := sqlcgen.ListSubscriptionRow{}
	copier.Copy(&row, &subscription)

	output := ToListSubscriptionOutput(row)
	return output, nil
}

func (s *SubscriptionService) GetSubscriptionByOrgId(ctx context.Context, orgId int, limit int, offset int) (ListSubscriptionOutputWithCount, error) {

	input := sqlcgen.GetSubscriptionByOrgIdParams{
		OrganizationID: int32(orgId),
		Limit:          int32(limit),
		Offset:         int32(offset),
	}

	subscriptions, err := s.DB.GetSubscriptionByOrgId(ctx, input)
	if err != nil {
		return ListSubscriptionOutputWithCount{}, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the subscriptions",
			err,
		)
	}

	listSubscriptionRows := make([]sqlcgen.ListSubscriptionRow, len(subscriptions))
	for i, sub := range subscriptions {
		listSubscriptionRows[i] = sqlcgen.ListSubscriptionRow(sub)
	}

	output := ToListSubscriptionOutputWithCount(listSubscriptionRows)
	return output, nil
}

func (s *SubscriptionService) ListSubscription(ctx context.Context, limit int, offset int) (ListSubscriptionOutputWithCount, error) {

	input := sqlcgen.ListSubscriptionParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	subscriptions, err := s.DB.ListSubscription(ctx, input)
	if err != nil {
		return ListSubscriptionOutputWithCount{}, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the subscriptions",
			err,
		)
	}

	output := ToListSubscriptionOutputWithCount(subscriptions)
	return output, nil
}

func (s *SubscriptionService) UpdateSubscription(ctx context.Context, input *UpdateSubscriptionParams) error {

	params := sqlcgen.UpdateSubscriptionParams{
		SubscriptionName:               input.Name,
		SubscriptionStartDate:          int32(input.StartDate.Unix()),
		SubscriptionApiLimit:           converter.ToPgInt4(input.APILimit),
		SubscriptionExpiryDate:         converter.ToPgInt4FromTimeOrDate(input.ExpiryDate),
		SubscriptionDescription:        converter.ToPgText(input.Description),
		SubscriptionStatus:             converter.ToPgBool(input.Status),
		OrganizationID:                 int32(input.OrganizationID),
		SubscriptionTierID:             int32(input.SubscriptionTierID),
		SubscriptionBillingInterval:    converter.ToPgText(input.BillingInterval),
		SubscriptionBillingModel:       converter.ToPgText(input.BillingModel),
		SubscriptionQuotaResetInterval: converter.ToPgText(input.QuotaResetInterval),
		SubscriptionID:                 int32(input.SubscriptionID),
	}

	_, err := s.DB.UpdateSubscription(ctx, params)
	if err != nil {
		return server.NewError(
			server.ErrorInternal,
			"couldn't update the organization",
			err,
		)
	}

	err = s.PubSubClient.Publish(ctx, string(common.SubscriptionModified), common.SubscriptionModifiedEvent{
		ID: int32(input.OrganizationID),
	})
	if err != nil {
		return server.NewError(
			server.ErrorInternal,
			"couldn't push to the queue",
			err,
		)
	}

	return nil
}
