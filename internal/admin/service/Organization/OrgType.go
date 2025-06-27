package organization

import (
	"context"

	"github.com/bignyap/go-admin/internal/database/dbutils"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/server"
	"github.com/jackc/pgx/v5"
)

type BulkCreateOrgTypeInserter struct {
	OrgTypes            CreateOrgTypeParams
	OrganizationService *OrganizationService
}

func (input BulkCreateOrgTypeInserter) InsertRows(ctx context.Context, tx pgx.Tx) (int64, error) {
	return input.OrganizationService.DB.CreateOrgTypes(ctx, input.OrgTypes.Names)
}

func (s *OrganizationService) CreateOrgTypeInBatch(ctx context.Context, input CreateOrgTypeParams) (int64, error) {

	inserter := BulkCreateOrgTypeInserter{
		OrgTypes:            input,
		OrganizationService: s,
	}

	affectedRows, err := dbutils.InsertWithTransaction(ctx, s.Conn, inserter)
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal,
			"couldn't create the organization types",
			err,
		)
	}

	return affectedRows, nil
}

func (s *OrganizationService) CreateOrgType(ctx context.Context, name string) (CreateOrgTypeOutput, error) {

	insertedID, err := s.DB.CreateOrgType(ctx, name)
	if err != nil {
		return CreateOrgTypeOutput{}, server.NewError(
			server.ErrorInternal,
			"couldn't create the organization type",
			err,
		)
	}

	output := CreateOrgTypeOutput{
		ID: int(insertedID),
		CreateOrgTypeInput: CreateOrgTypeInput{
			Name: name,
		},
	}

	return output, nil
}

func (s *OrganizationService) ListOrgType(ctx context.Context, limit int, offset int) ([]CreateOrgTypeOutput, error) {

	input := sqlcgen.ListOrgTypeParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	orgTypes, err := s.DB.ListOrgType(ctx, input)
	if err != nil {
		return nil, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the organization types",
			err,
		)
	}

	if len(orgTypes) == 0 {
		return []CreateOrgTypeOutput{}, nil
	}

	var output []CreateOrgTypeOutput
	for _, orgType := range orgTypes {
		output = append(output, CreateOrgTypeOutput{
			ID: int(orgType.OrganizationTypeID),
			CreateOrgTypeInput: CreateOrgTypeInput{
				Name: orgType.OrganizationTypeName,
			},
		})
	}

	return output, nil
}

func (s *OrganizationService) DeleteOrgType(ctx context.Context, typeId int) error {

	if err := s.DB.DeleteOrgTypeById(ctx, int32(typeId)); err != nil {
		return server.NewError(
			server.ErrorInternal,
			"couldn't delete the organization type",
			err,
		)
	}

	return nil
}
