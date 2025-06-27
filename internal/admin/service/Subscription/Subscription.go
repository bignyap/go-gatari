package subscription

import (
	"context"
	"time"

	"github.com/bignyap/go-admin/internal/database/dbutils"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/server"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
		SubscriptionName:        input.Name,
		SubscriptionType:        input.Type,
		SubscriptionCreatedDate: int32(input.CreatedAt.Unix()),
		SubscriptionUpdatedDate: int32(input.UpdatedAt.Unix()),
		SubscriptionStartDate:   int32(input.StartDate.Unix()),
		SubscriptionApiLimit:    toPgInt4(input.APILimit),
		SubscriptionExpiryDate:  toPgInt4FromTimePtr(input.ExpiryDate),
		SubscriptionDescription: toPgText(input.Description),
		SubscriptionStatus:      toPgBool(input.Status),
		OrganizationID:          int32(input.OrganizationID),
		SubscriptionTierID:      int32(input.SubscriptionTierID),
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
			SubscriptionName:        input.Name,
			SubscriptionType:        input.Type,
			SubscriptionCreatedDate: currentTime,
			SubscriptionUpdatedDate: currentTime,
			SubscriptionStartDate:   currentTime,
			SubscriptionApiLimit:    toPgInt4(input.APILimit),
			SubscriptionExpiryDate:  toPgInt4FromTimePtr(input.ExpiryDate),
			SubscriptionDescription: toPgText(input.Description),
			SubscriptionStatus:      toPgBool(input.Status),
			OrganizationID:          int32(input.OrganizationID),
			SubscriptionTierID:      int32(input.SubscriptionTierID),
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

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, subId int) error {

	err := s.DB.DeleteOrganizationById(ctx, int32(subId))
	if err != nil {
		return server.NewError(
			server.ErrorInternal,
			"couldn't delete the subscription",
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

// -------- pgtype helpers --------

func toPgInt4(ptr *int) pgtype.Int4 {
	if ptr == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*ptr), Valid: true}
}

// func toPgInt4Ptr(value int) pgtype.Int4 {
// 	return pgtype.Int4{Int32: int32(value), Valid: true}
// }

func toPgInt4Ptr(v *int) pgtype.Int4 {
	if v == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*v), Valid: true}
}

func toPgInt4FromTime(t time.Time) pgtype.Int4 {
	return pgtype.Int4{Int32: int32(t.Unix()), Valid: true}
}

func toPgInt4FromTimePtr(ptr *time.Time) pgtype.Int4 {
	if ptr == nil {
		return pgtype.Int4{Valid: false}
	}
	return toPgInt4FromTime(*ptr)
}

func toPgText(ptr *string) pgtype.Text {
	if ptr == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *ptr, Valid: true}
}

func toPgBool(ptr *bool) pgtype.Bool {
	if ptr == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *ptr, Valid: true}
}

func fromPgInt4Ptr(v pgtype.Int4) *int {
	if !v.Valid {
		return nil
	}
	val := int(v.Int32)
	return &val
}

func fromPgInt4TimePtr(v pgtype.Int4) *time.Time {
	if !v.Valid {
		return nil
	}
	t := time.Unix(int64(v.Int32), 0)
	return &t
}

func fromPgText(v pgtype.Text) *string {
	if !v.Valid {
		return nil
	}
	return &v.String
}

func fromPgBool(v pgtype.Bool) *bool {
	if !v.Valid {
		return nil
	}
	return &v.Bool
}
