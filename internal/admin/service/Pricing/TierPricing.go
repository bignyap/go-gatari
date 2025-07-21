package pricing

import (
	"context"
	"strings"

	"github.com/bignyap/go-admin/internal/common"
	"github.com/bignyap/go-admin/internal/database/dbutils"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/converter"
	"github.com/bignyap/go-utilities/server"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type BulkCreateTierPricingsInserter struct {
	TierPricings   []sqlcgen.CreateTierPricingsParams
	PricingService *PricingService
}

func (input BulkCreateTierPricingsInserter) InsertRows(ctx context.Context, tx pgx.Tx) (int64, error) {
	return input.PricingService.DB.CreateTierPricings(ctx, input.TierPricings)
}

func (s *PricingService) CreateTierPricingInBatch(ctx context.Context, input []sqlcgen.CreateTierPricingsParams) (int, error) {

	inserter := BulkCreateTierPricingsInserter{
		TierPricings:   input,
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

func (s *PricingService) CreateTierPricing(ctx context.Context, input *sqlcgen.CreateTierPricingParams) (CreateTierPricingOutput, error) {

	insertedID, err := s.DB.CreateTierPricing(ctx, *input)
	if err != nil {
		return CreateTierPricingOutput{}, server.NewError(
			server.ErrorInternal,
			"couldn't create the custom pricing",
			err,
		)
	}

	output := CreateTierPricingOutput{
		ID: int(insertedID),
		CreateTierPricingParams: CreateTierPricingParams{
			BaseCostPerCall:    input.BaseCostPerCall,
			BaseRateLimit:      converter.FromPgInt4Ptr(input.BaseRateLimit),
			ApiEndpointId:      int(input.ApiEndpointID),
			SubscriptionTierID: int(input.SubscriptionTierID),
			CostMode:           input.CostMode,
		},
	}

	return output, nil
}

func (s *PricingService) GetTierPricingByTierId(ctx context.Context, id int, limit int, offset int) (CreateTierPricingOutputWithCount, error) {

	input := sqlcgen.GetTierPricingByTierIdParams{
		SubscriptionTierID: int32(id),
		Limit:              int32(limit),
		Offset:             int32(offset),
	}

	tierPricings, err := s.DB.GetTierPricingByTierId(ctx, input)
	if err != nil {
		return CreateTierPricingOutputWithCount{}, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the tier pricing list",
			err,
		)
	}

	output := make([]CreateTierPricingWithTierName, len(tierPricings))
	for i, tierPricing := range tierPricings {
		output[i] = CreateTierPricingWithTierName{
			EndpointName: tierPricing.EndpointName,
			CreateTierPricingOutput: CreateTierPricingOutput{
				ID: int(tierPricing.TierBasePricingID),
				CreateTierPricingParams: CreateTierPricingParams{
					SubscriptionTierID: int(tierPricing.SubscriptionTierID),
					ApiEndpointId:      int(tierPricing.ApiEndpointID),
					BaseCostPerCall:    tierPricing.BaseCostPerCall,
					BaseRateLimit:      fromPgInt4Ptr(tierPricing.BaseRateLimit),
					CostMode:           tierPricing.CostMode,
				},
			},
		}
	}

	totalItems := 0
	if len(tierPricings) > 0 {
		totalItems = int(tierPricings[0].TotalItems)
	}

	return CreateTierPricingOutputWithCount{
		Data:       output,
		TotalItems: totalItems,
	}, nil
}

func fromPgInt4Ptr(v pgtype.Int4) *int {
	if !v.Valid {
		return nil
	}
	val := int(v.Int32)
	return &val
}

func (s *PricingService) DeleteTierPricing(ctx context.Context, idType string, id int) error {

	var pubsubId int32

	switch strings.ToLower(idType) {
	case "id":
		tier, err := s.DB.DeleteTierPricingById(ctx, int32(id))
		if err != nil {
			return server.NewError(
				server.ErrorInternal,
				"couldn't delete the subscription by organization_id",
				err,
			)
		}
		pubsubId = tier
	case "tier":
		tier, err := s.DB.DeleteTierPricingByTierId(ctx, int32(id))
		if err != nil {
			return server.NewError(
				server.ErrorInternal,
				"couldn't delete the subscription by id",
				err,
			)
		}
		pubsubId = tier
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
