package billing

import (
	"context"
	"time"

	"github.com/bignyap/go-admin/database/dbutils"
	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-utilities/server"
	"github.com/jackc/pgx/v5"
)

type BulkBillingHistoryInserter struct {
	BillingHistories []sqlcgen.CreateBillingHistoriesParams
	BillingService   *BillingService
}

func (input BulkBillingHistoryInserter) InsertRows(ctx context.Context, tx pgx.Tx) (int64, error) {
	return input.BillingService.DB.CreateBillingHistories(ctx, input.BillingHistories)
}

func (s *BillingService) CreateBillingHistory(ctx context.Context, input CreateBillingHistoryParams) (CreateBillingHistoryOutput, error) {

	input.CreatedAt = time.Now()

	if err := s.Validator.Struct(input); err != nil {
		return CreateBillingHistoryOutput{}, server.NewError(
			server.ErrorBadRequest,
			"validation error",
			err,
		)
	}

	sqlInput := sqlcgen.CreateBillingHistoryParams{
		BillingStartDate: int32(input.StartDate.Unix()),
		BillingEndDate:   int32(input.EndDate.Unix()),
		TotalAmountDue:   input.TotalAmountDue,
		TotalCalls:       int32(input.TotalCalls),
		PaymentStatus:    input.PaymentStatus,
		PaymentDate:      timePtrToInt4(input.PaymentDate),
		BillingCreatedAt: int32(input.CreatedAt.Unix()),
		SubscriptionID:   int32(input.SubscriptionId),
	}

	billingID, err := s.DB.CreateBillingHistory(ctx, sqlInput)
	if err != nil {
		return CreateBillingHistoryOutput{}, server.NewError(
			server.ErrorInternal,
			"couldn't create the billing history",
			err,
		)
	}

	output := CreateBillingHistoryOutput{
		ID:                         int(billingID),
		CreateBillingHistoryParams: input,
	}

	return output, nil
}

func (s *BillingService) CreateBillingHistoryInBatch(ctx context.Context, input []sqlcgen.CreateBillingHistoriesParams) (int, error) {

	inserter := BulkBillingHistoryInserter{
		BillingHistories: input,
		BillingService:   s,
	}

	affectedRows, err := dbutils.InsertWithTransaction(ctx, s.Conn, inserter)
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal,
			"couldn't create the billing histories",
			err,
		)
	}

	return int(affectedRows), nil
}

func (s *BillingService) GetBillingHistoryByOrgId(ctx context.Context, orgId int, n int, page int) ([]sqlcgen.BillingHistory, error) {

	input := sqlcgen.GetBillingHistoryByOrgIdParams{
		OrganizationID: int32(orgId),
		Limit:          int32(page),
		Offset:         int32(n),
	}

	billingHistories, err := s.DB.GetBillingHistoryByOrgId(ctx, input)
	if err != nil {
		return nil, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the billing histories",
			err,
		)
	}

	return billingHistories, nil
}

func (s *BillingService) GetBillingHistoryBySubId(ctx context.Context, subId int, n int, page int) ([]sqlcgen.BillingHistory, error) {

	input := sqlcgen.GetBillingHistoryBySubIdParams{
		SubscriptionID: int32(subId),
		Limit:          int32(page),
		Offset:         int32(n),
	}

	billingHistories, err := s.DB.GetBillingHistoryBySubId(ctx, input)
	if err != nil {
		return nil, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the billing histories",
			err,
		)
	}

	return billingHistories, nil
}

func (s *BillingService) GetBillingHistoryById(ctx context.Context, id int, n int, page int) ([]sqlcgen.BillingHistory, error) {

	input := sqlcgen.GetBillingHistoryByIdParams{
		BillingID: int32(id),
		Limit:     int32(page),
		Offset:    int32(n),
	}

	billingHistories, err := s.DB.GetBillingHistoryById(ctx, input)
	if err != nil {
		return nil, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the billing histories",
			err,
		)
	}

	return billingHistories, nil
}
