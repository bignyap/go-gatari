package organization

import (
	"context"
	"fmt"

	"github.com/bignyap/go-admin/database/dbutils"
	"github.com/bignyap/go-admin/database/sqlcgen"
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
		return 0, fmt.Errorf("couldn't create the organization types: %s", err)
	}

	return affectedRows, nil
}

func (s *OrganizationService) CreateOrgType(ctx context.Context, name string) (CreateOrgTypeOutput, error) {

	insertedID, err := s.DB.CreateOrgType(ctx, name)
	if err != nil {
		return CreateOrgTypeOutput{}, fmt.Errorf("couldn't create the organization type: %s", err)
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
		return nil, fmt.Errorf("couldn't retrieve the organization types: %s", err)
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
		return fmt.Errorf("couldn't delete the organization type: %s", err)
	}

	return nil
}
