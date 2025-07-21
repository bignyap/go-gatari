package organization

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/bignyap/go-admin/internal/common"
	"github.com/bignyap/go-admin/internal/database/dbutils"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/converter"

	"github.com/bignyap/go-utilities/server"
)

func (apiCfg *OrganizationService) CreateOrganization(ctx context.Context, input *CreateOrganizationParams) (CreateOrganizationOutput, error) {

	currentTime := int32(converter.ToUnixTime())
	org := sqlcgen.CreateOrganizationParams{
		OrganizationName:         input.Name,
		OrganizationRealm:        input.Realm,
		OrganizationSupportEmail: input.SupportEmail,
		OrganizationCreatedAt:    currentTime,
		OrganizationUpdatedAt:    currentTime,
		OrganizationCountry:      converter.ToPgText(input.Country),
		OrganizationConfig:       converter.ToPgText(input.Config),
		OrganizationActive:       converter.ToPgBool(input.Active),
		OrganizationReportQ:      converter.ToPgBool(input.ReportQ),
		OrganizationTypeID:       int32(input.TypeID),
	}

	insertedID, err := apiCfg.DB.CreateOrganization(ctx, org)
	if err != nil {
		return CreateOrganizationOutput{}, server.NewError(
			server.ErrorInternal,
			"couldn't create the organization",
			err,
		)
	}

	input.CreatedAt = converter.FromUnixTime32(currentTime)
	input.UpdatedAt = converter.FromUnixTime32(currentTime)

	return CreateOrganizationOutput{
		ID:                       int(insertedID),
		CreateOrganizationParams: *input,
	}, nil
}

func (apiCfg *OrganizationService) CreateOrganizationInBatch(ctx context.Context, inputs []CreateOrganizationParams) (int, error) {

	currentTime := int32(converter.ToUnixTime())
	var batch []sqlcgen.CreateOrganizationsParams

	for _, input := range inputs {
		batch = append(batch, sqlcgen.CreateOrganizationsParams{
			OrganizationName:         input.Name,
			OrganizationRealm:        input.Realm,
			OrganizationSupportEmail: input.SupportEmail,
			OrganizationCreatedAt:    currentTime,
			OrganizationUpdatedAt:    currentTime,
			OrganizationCountry:      converter.ToPgText(input.Country),
			OrganizationConfig:       converter.ToPgText(input.Config),
			OrganizationActive:       converter.ToPgBool(input.Active),
			OrganizationReportQ:      converter.ToPgBool(input.ReportQ),
			OrganizationTypeID:       int32(input.TypeID),
		})
	}

	inserter := BulkOrganizationInserter{
		Organizations:       batch,
		OrganizationService: apiCfg,
	}

	affectedRows, err := dbutils.InsertWithTransaction(ctx, apiCfg.Conn, inserter)
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal,
			"couldn't create the organizations",
			err,
		)
	}

	return int(affectedRows), nil
}

type BulkOrganizationInserter struct {
	Organizations       []sqlcgen.CreateOrganizationsParams
	OrganizationService *OrganizationService
}

func (input BulkOrganizationInserter) InsertRows(ctx context.Context, tx pgx.Tx) (int64, error) {
	return input.OrganizationService.DB.CreateOrganizations(ctx, input.Organizations)
}

func (s *OrganizationService) ListOrganizations(ctx context.Context, limit int, offset int) (ListOrganizationOutputWithCount, error) {

	input := sqlcgen.ListOrganizationParams{
		Limit:   int32(limit),
		Offset:  int32(offset),
		Column1: 0, // it's equivalent to 0. Workaround for sqlc issue
	}

	organizations, err := s.DB.ListOrganization(ctx, input)
	if err != nil {
		return ListOrganizationOutputWithCount{}, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the organizations",
			err,
		)
	}

	output := ToListOrganizationOutputWithCount(organizations)
	return output, nil
}

func (s *OrganizationService) GetOrganizationById(ctx context.Context, orgId int) (ListOrganizationOutput, error) {

	input := sqlcgen.ListOrganizationParams{
		Limit:          1,
		Offset:         0,
		OrganizationID: int32(orgId),
		Column1:        int32(orgId),
	}

	organization, err := s.DB.ListOrganization(ctx, input)
	if err != nil {
		return ListOrganizationOutput{}, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the organizations",
			err,
		)
	}

	if len(organization) == 0 {
		return ListOrganizationOutput{}, nil
	}

	output := ToListOrganizationOutput(organization[0])
	return output, nil
}

func (s *OrganizationService) DeleteOrganizationById(ctx context.Context, id int) error {

	err := s.DB.DeleteOrganizationById(ctx, int32(id))
	if err != nil {
		return server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the organizations",
			err,
		)
	}

	return nil
}

func (apiCfg *OrganizationService) UpdateOrganization(ctx context.Context, input *UpdateOrganizationParams) error {

	currentTime := int32(converter.ToUnixTime())
	org := sqlcgen.UpdateOrganizationParams{

		OrganizationRealm:        input.Realm,
		OrganizationName:         input.Name,
		OrganizationSupportEmail: input.SupportEmail,
		OrganizationUpdatedAt:    currentTime,
		OrganizationCountry:      converter.ToPgText(input.Country),
		OrganizationConfig:       converter.ToPgText(input.Config),
		OrganizationActive:       converter.ToPgBool(input.Active),
		OrganizationReportQ:      converter.ToPgBool(input.ReportQ),
		OrganizationTypeID:       int32(input.TypeID),
		OrganizationID:           int32(input.OrganizationID),
	}

	_, err := apiCfg.DB.UpdateOrganization(ctx, org)
	if err != nil {
		return server.NewError(
			server.ErrorInternal,
			"couldn't update the organization",
			err,
		)
	}

	err = apiCfg.PubSubClient.Publish(ctx, string(common.OrganizationModified), common.OrganizationModifiedEvent{
		ID:   int32(input.OrganizationID),
		Name: input.Realm,
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
