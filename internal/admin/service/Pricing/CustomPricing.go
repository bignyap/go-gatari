package pricing

import (
	"context"
	"strings"

	"github.com/bignyap/go-admin/internal/common"
	"github.com/bignyap/go-admin/internal/database/dbutils"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/server"
	"github.com/jackc/pgx/v5"
)

type BulkCreateCustomPricingsInserter struct {
	CustomPricings []sqlcgen.CreateCustomPricingsParams
	PricingService *PricingService
}

func (input BulkCreateCustomPricingsInserter) InsertRows(ctx context.Context, tx pgx.Tx) (int64, error) {

	affectedRows, err := input.PricingService.DB.CreateCustomPricings(ctx, input.CustomPricings)
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (s *PricingService) CreateCustomPricingInBatch(ctx context.Context, input []sqlcgen.CreateCustomPricingsParams) (int, error) {

	inserter := BulkCreateCustomPricingsInserter{
		CustomPricings: input,
		PricingService: s,
	}

	affectedRows, err := dbutils.InsertWithTransaction(ctx, s.Conn, inserter)
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal,
			"couldn't create the tier pricings",
			err,
		)
	}

	return int(affectedRows), nil
}

func (s *PricingService) CreateCustomPricing(ctx context.Context, input *sqlcgen.CreateCustomPricingParams) (CreateCustomPricingOutput, error) {

	insertedID, err := s.DB.CreateCustomPricing(ctx, *input)
	if err != nil {
		return CreateCustomPricingOutput{}, server.NewError(
			server.ErrorInternal,
			"couldn't create the custom pricing",
			err,
		)
	}

	output := CreateCustomPricingOutput{
		ID: int(insertedID),
		CreateCustomPricingParams: CreateCustomPricingParams{
			CustomCostPerCall: input.CustomCostPerCall,
			CustomRateLimit:   int(input.CustomRateLimit),
			SubscriptionID:    int(input.SubscriptionID),
			TierBasePricingID: int(input.TierBasePricingID),
			CostMode:          input.CostMode,
		},
	}

	return output, nil
}

func (s *PricingService) DeleteCustomPricing(ctx context.Context, idType string, id int) error {

	var pubsubId int32

	switch strings.ToLower(idType) {
	case "subscription":
		err := s.DB.DeleteCustomPricingBySubscriptionId(ctx, int32(id))
		if err != nil {
			return server.NewError(
				server.ErrorInternal,
				"couldn't delete the custom pricing by subscription_id",
				err,
			)
		}
	case "pricing":
		err := s.DB.DeleteCustomPricingById(ctx, int32(id))
		if err != nil {
			return server.NewError(
				server.ErrorInternal,
				"couldn't delete the custom pricing by id",
				err,
			)
		}
	}

	err := s.PubSubClient.Publish(ctx, string(common.SubscriptionModified), common.PricingModifiedEvent{
		ID:   pubsubId,
		Type: "subscription_tier",
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

func (s *PricingService) GetCustomPricing(ctx context.Context, sId int, limit int, offset int) ([]CreateCustomPricingOutput, error) {

	input := sqlcgen.GetCustomPricingParams{
		SubscriptionID: int32(sId),
		Limit:          int32(limit),
		Offset:         int32(offset),
	}

	customPricings, err := s.DB.GetCustomPricing(ctx, input)
	if err != nil {
		return []CreateCustomPricingOutput{}, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the custom pricing list",
			err,
		)
	}

	var output []CreateCustomPricingOutput

	for _, customPricing := range customPricings {

		output = append(output, CreateCustomPricingOutput{
			ID: int(customPricing.TierBasePricingID),
			CreateCustomPricingParams: CreateCustomPricingParams{
				TierBasePricingID: int(customPricing.TierBasePricingID),
				SubscriptionID:    int(customPricing.SubscriptionID),
				CustomCostPerCall: customPricing.CustomCostPerCall,
				CustomRateLimit:   int(customPricing.CustomRateLimit),
				CostMode:          customPricing.CostMode,
			},
		})
	}

	return output, nil
}
