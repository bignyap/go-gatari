package subscription

import (
	"context"
	"time"

	"github.com/bignyap/go-admin/internal/database/dbutils"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/server"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type BulkCreateSubscriptionTierInserter struct {
	SubscriptionTiers   []sqlcgen.CreateSubscriptionTiersParams
	SubscriptionService *SubscriptionService
}

func (input BulkCreateSubscriptionTierInserter) InsertRows(ctx context.Context, tx pgx.Tx) (int64, error) {
	return input.SubscriptionService.DB.CreateSubscriptionTiers(ctx, input.SubscriptionTiers)
}

func (s *SubscriptionService) CreateSubscriptionTierInBatch(ctx context.Context, input []sqlcgen.CreateSubscriptionTiersParams) (int, error) {

	inserter := BulkCreateSubscriptionTierInserter{
		SubscriptionTiers:   input,
		SubscriptionService: s,
	}

	affectedRows, err := dbutils.InsertWithTransaction(ctx, s.Conn, inserter)
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal,
			"couldn't create the subscription tiers",
			err,
		)
	}

	return int(affectedRows), nil
}

func (s *SubscriptionService) CreateSubscriptionTier(ctx context.Context, input CreateSubTierParams) (CreateSubTierOutput, error) {

	if err := s.Validator.Struct(input); err != nil {
		return CreateSubTierOutput{}, server.NewError(
			server.ErrorInternal,
			"validation failed",
			err,
		)
	}

	currentTime := int32(time.Now().Unix())
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	var description pgtype.Text
	if input.Description != nil {
		description.String = *input.Description
		description.Valid = true
	}

	subTierParams := sqlcgen.CreateSubscriptionTierParams{
		TierName:        input.Name,
		TierDescription: description,
		TierCreatedAt:   currentTime,
		TierUpdatedAt:   currentTime,
	}

	err := s.DB.ArchiveExistingSubscriptionTier(ctx, input.Name)
	if err != nil {
		return CreateSubTierOutput{}, server.NewError(
			server.ErrorInternal,
			"couldn't archive the existing subscription tier",
			err,
		)
	}

	insertedID, err := s.DB.CreateSubscriptionTier(ctx, subTierParams)
	if err != nil {
		return CreateSubTierOutput{}, server.NewError(
			server.ErrorInternal,
			"couldn't create the subscription tier",
			err,
		)
	}

	output := CreateSubTierOutput{
		ID:       int(insertedID),
		Archived: false,
		CreateSubTierParams: CreateSubTierParams{
			Name:        input.Name,
			Description: input.Description,
			CreatedAt:   input.CreatedAt,
			UpdatedAt:   input.UpdatedAt,
		},
	}
	return output, nil
}

func (s *SubscriptionService) ListSubscriptionTiers(ctx context.Context, limit int, offset int, archived bool) (CreateSubTierOutputWithCount, error) {

	input := sqlcgen.ListSubscriptionTierParams{
		TierArchived: archived,
		Limit:        int32(limit),
		Offset:       int32(offset),
	}

	subTiers, err := s.DB.ListSubscriptionTier(ctx, input)
	if err != nil {
		return CreateSubTierOutputWithCount{}, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the subscription tiers",
			err,
		)
	}

	var output []CreateSubTierOutput
	for _, subTier := range subTiers {
		var description *string
		if subTier.TierDescription.Valid {
			description = &subTier.TierDescription.String
		}

		output = append(output, CreateSubTierOutput{
			ID:       int(subTier.SubscriptionTierID),
			Archived: subTier.TierArchived,
			CreateSubTierParams: CreateSubTierParams{
				Name:        subTier.TierName,
				Description: description,
				CreatedAt:   time.Unix(int64(subTier.TierCreatedAt), 0),
				UpdatedAt:   time.Unix(int64(subTier.TierUpdatedAt), 0),
			},
		})
	}

	totalItems := 0
	if len(subTiers) > 0 {
		totalItems = int(subTiers[0].TotalItems)
	}

	return CreateSubTierOutputWithCount{
		Data:       output,
		TotalItems: totalItems,
	}, nil
}

func (s *SubscriptionService) DeleteSubscriptionTier(ctx context.Context, id int) error {

	err := s.DB.DeleteSubscriptionTierById(ctx, int32(id))
	if err != nil {
		return server.NewError(
			server.ErrorInternal,
			"couldn't delete the subscription tier",
			err,
		)
	}

	return nil
}
