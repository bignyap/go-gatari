package usage

import (
	"github.com/bignyap/go-admin/internal/database/sqlcgen"

	"context"

	"github.com/bignyap/go-admin/internal/database/dbutils"
	"github.com/bignyap/go-utilities/server"
	"github.com/jackc/pgx/v5"
)

type BulkApiSummaryInserter struct {
	ApiUsageSummaries   []sqlcgen.CreateApiUsageSummariesParams
	UsageSummaryService *UsageSummaryService
}

func (input BulkApiSummaryInserter) InsertRows(ctx context.Context, tx pgx.Tx) (int64, error) {

	affectedRows, err := input.UsageSummaryService.DB.CreateApiUsageSummaries(ctx, input.ApiUsageSummaries)
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (s *UsageSummaryService) CreateApiUsageInBatch(ctx context.Context, input []sqlcgen.CreateApiUsageSummariesParams) (int64, error) {

	inserter := BulkApiSummaryInserter{
		ApiUsageSummaries:   input,
		UsageSummaryService: s,
	}

	affectedRows, err := dbutils.InsertWithTransaction(ctx, s.Conn, inserter)
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal,
			"couldn't create the api usage summaries",
			err,
		)
	}
	return affectedRows, nil
}

func (s *UsageSummaryService) CreateApiUsage(ctx context.Context, input *sqlcgen.CreateApiUsageSummaryParams) (int32, error) {

	insertedID, err := s.DB.CreateApiUsageSummary(ctx, *input)
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal,
			"couldn't create the API usage summary",
			err,
		)
	}

	return insertedID, nil
}

func (s *UsageSummaryService) GetApiUsageSummaryByOrgId(ctx context.Context, orgId int, n int, page int) ([]CreateApiUsageSummaryOutput, error) {

	input := sqlcgen.GetApiUsageSummaryByOrgIdParams{
		OrganizationID: int32(orgId),
		Limit:          int32(page),
		Offset:         int32(n),
	}

	apiUsageSummaries, err := s.DB.GetApiUsageSummaryByOrgId(ctx, input)
	if err != nil {
		return nil, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the api usage summaries",
			err,
		)
	}

	var output []CreateApiUsageSummaryOutput

	for _, apiUsageSummary := range apiUsageSummaries {

		output = append(output, ToCreateApiUsageSummaryOutput(apiUsageSummary))
	}

	return output, nil
}

func (s *UsageSummaryService) GetApiUsageSummaryBySubId(ctx context.Context, subId int, n int, page int) ([]CreateApiUsageSummaryOutput, error) {

	input := sqlcgen.GetApiUsageSummaryBySubIdParams{
		SubscriptionID: int32(subId),
		Limit:          int32(page),
		Offset:         int32(n),
	}

	apiUsageSummaries, err := s.DB.GetApiUsageSummaryBySubId(ctx, input)
	if err != nil {
		return nil, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the api usage summaries",
			err,
		)
	}

	var output []CreateApiUsageSummaryOutput

	for _, apiUsageSummary := range apiUsageSummaries {

		output = append(output, ToCreateApiUsageSummaryOutput(apiUsageSummary))
	}

	return output, nil
}

func (s *UsageSummaryService) GetApiUsageSummaryByEndpointId(ctx context.Context, eId int, n int, page int) ([]CreateApiUsageSummaryOutput, error) {

	input := sqlcgen.GetApiUsageSummaryByEndpointIdParams{
		ApiEndpointID: int32(eId),
		Limit:         int32(page),
		Offset:        int32(n),
	}

	apiUsageSummaries, err := s.DB.GetApiUsageSummaryByEndpointId(ctx, input)
	if err != nil {
		return nil, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the api usage summaries",
			err,
		)
	}

	var output []CreateApiUsageSummaryOutput

	for _, apiUsageSummary := range apiUsageSummaries {

		output = append(output, ToCreateApiUsageSummaryOutput(apiUsageSummary))
	}

	return output, nil
}
