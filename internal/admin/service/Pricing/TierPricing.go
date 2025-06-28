package pricing

import (
	"context"
	"strings"

	"github.com/bignyap/go-admin/internal/database/dbutils"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
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

	switch strings.ToLower(idType) {
	case "id":
		err := s.DB.DeleteTierPricingById(ctx, int32(id))
		if err != nil {
			return server.NewError(
				server.ErrorInternal,
				"couldn't delete the subscription by organization_id",
				err,
			)
		}
	case "tier":
		err := s.DB.DeleteTierPricingByTierId(ctx, int32(id))
		if err != nil {
			return server.NewError(
				server.ErrorInternal,
				"couldn't delete the subscription by id",
				err,
			)
		}
	}

	return nil
}
