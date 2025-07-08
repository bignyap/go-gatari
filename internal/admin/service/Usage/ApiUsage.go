package usage

import (
	"context"

	"github.com/bignyap/go-admin/internal/database/dbutils"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/converter"
	"github.com/bignyap/go-utilities/server"
	"github.com/jackc/pgx/v5"
)

type BulkApiSummaryInserter struct {
	ApiUsageSummaries   []sqlcgen.CreateApiUsageSummariesParams
	UsageSummaryService *UsageSummaryService
}

func (input BulkApiSummaryInserter) InsertRows(ctx context.Context, tx pgx.Tx) (int64, error) {
	return input.UsageSummaryService.DB.CreateApiUsageSummaries(ctx, input.ApiUsageSummaries)
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

func (s *UsageSummaryService) GetUsageSummary(ctx context.Context, filters UsageSummaryFilters) ([]sqlcgen.GetUsageSummaryRow, error) {
	return GetGroupedUsageSummary(ctx, s.DB.GetUsageSummary, buildGetUsageSummaryInputParams, filters)
}

func (s *UsageSummaryService) GetUsageSummaryByDay(ctx context.Context, filters UsageSummaryFilters) ([]sqlcgen.GetUsageSummaryGroupedByDayRow, error) {
	return GetGroupedUsageSummary(ctx, s.DB.GetUsageSummaryGroupedByDay, buildGetUsageSummaryByDayInputParams, filters)
}

func buildGetUsageSummaryInputParams(f UsageSummaryFilters) sqlcgen.GetUsageSummaryParams {
	val := sqlcgen.GetUsageSummaryParams{
		OrgID:      converter.ToPgInt4(f.OrgID),
		SubID:      converter.ToPgInt4(f.SubID),
		EndpointID: converter.ToPgInt4(f.EndpointID),
		StartDate:  converter.ToPgInt4(f.StartDate),
		EndDate:    converter.ToPgInt4(f.EndDate),
		Limit:      int32(f.Limit),
		Offset:     int32(f.Offset),
	}
	return val
}

func buildGetUsageSummaryByDayInputParams(f UsageSummaryFilters) sqlcgen.GetUsageSummaryGroupedByDayParams {
	return sqlcgen.GetUsageSummaryGroupedByDayParams{
		OrgID:      converter.ToPgInt4(f.OrgID),
		SubID:      converter.ToPgInt4(f.SubID),
		EndpointID: converter.ToPgInt4(f.EndpointID),
		StartDate:  converter.ToPgInt4(f.StartDate),
		EndDate:    converter.ToPgInt4(f.EndDate),
		Limit:      int32(f.Limit),
		Offset:     int32(f.Offset),
	}
}
