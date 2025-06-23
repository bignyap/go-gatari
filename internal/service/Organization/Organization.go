package organization

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/bignyap/go-admin/database/dbutils"
	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-admin/utils/misc"
)

func (apiCfg *OrganizationService) CreateOrganization(ctx context.Context, input *CreateOrganizationParams) (CreateOrganizationOutput, error) {

	currentTime := int32(misc.ToUnixTime())
	org := sqlcgen.CreateOrganizationParams{
		OrganizationName:         input.Name,
		OrganizationRealm:        input.Realm,
		OrganizationSupportEmail: input.SupportEmail,
		OrganizationCreatedAt:    currentTime,
		OrganizationUpdatedAt:    currentTime,
		OrganizationCountry:      toText(input.Country),
		OrganizationConfig:       toText(input.Config),
		OrganizationActive:       toBool(input.Active),
		OrganizationReportQ:      toBool(input.ReportQ),
		OrganizationTypeID:       int32(input.TypeID),
	}

	insertedID, err := apiCfg.DB.CreateOrganization(ctx, org)
	if err != nil {
		return CreateOrganizationOutput{}, fmt.Errorf("couldn't create the organization: %s", err)
	}

	input.CreatedAt = misc.FromUnixTime32(currentTime)
	input.UpdatedAt = misc.FromUnixTime32(currentTime)

	return CreateOrganizationOutput{
		ID:                       int(insertedID),
		CreateOrganizationParams: *input,
	}, nil
}

func (apiCfg *OrganizationService) CreateOrganizationInBatch(ctx context.Context, inputs []CreateOrganizationParams) (int, error) {

	currentTime := int32(misc.ToUnixTime())
	var batch []sqlcgen.CreateOrganizationsParams

	for _, input := range inputs {
		batch = append(batch, sqlcgen.CreateOrganizationsParams{
			OrganizationName:         input.Name,
			OrganizationRealm:        input.Realm,
			OrganizationSupportEmail: input.SupportEmail,
			OrganizationCreatedAt:    currentTime,
			OrganizationUpdatedAt:    currentTime,
			OrganizationCountry:      toText(input.Country),
			OrganizationConfig:       toText(input.Config),
			OrganizationActive:       toBool(input.Active),
			OrganizationReportQ:      toBool(input.ReportQ),
			OrganizationTypeID:       int32(input.TypeID),
		})
	}

	inserter := BulkOrganizationInserter{
		Organizations:       batch,
		OrganizationService: apiCfg,
	}

	affectedRows, err := dbutils.InsertWithTransaction(ctx, apiCfg.Conn, inserter)
	if err != nil {
		return 0, fmt.Errorf("couldn't create the organizations: %s", err)
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
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	organizations, err := s.DB.ListOrganization(ctx, input)
	if err != nil {
		return ListOrganizationOutputWithCount{}, fmt.Errorf("couldn't retrieve the organizations: %s", err)
	}

	output := ToListOrganizationOutputWithCount(organizations)
	return output, nil
}

func (s *OrganizationService) GetOrganizationById(ctx context.Context, orgId int) (ListOrganizationOutput, error) {

	input := sqlcgen.ListOrganizationParams{
		Limit:          1,
		Offset:         0,
		OrganizationID: int32(orgId),
	}

	organization, err := s.DB.ListOrganization(ctx, input)
	if err != nil {
		return ListOrganizationOutput{}, fmt.Errorf("couldn't retrieve the organization: %s", err)
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
		return fmt.Errorf("couldn't delete the organization: %s", err)
	}

	return nil
}

func toText(val *string) pgtype.Text {
	if val == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *val, Valid: true}
}

func toBool(val *bool) pgtype.Bool {
	if val == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *val, Valid: true}
}
