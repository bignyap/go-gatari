package organization

import (
	"context"
	"fmt"
	"strings"

	"github.com/bignyap/go-admin/internal/database/dbutils"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/server"
	"github.com/jackc/pgx/v5"
)

type BulkCreateOrgPermissionInserter struct {
	OrgPermissions      []sqlcgen.CreateOrgPermissionsParams
	OrganizationService *OrganizationService
}

func (input BulkCreateOrgPermissionInserter) InsertRows(ctx context.Context, tx pgx.Tx) (int64, error) {
	return input.OrganizationService.DB.CreateOrgPermissions(ctx, input.OrgPermissions)
}

func (s *OrganizationService) CreateOrgPermissionInBatch(ctx context.Context, input []sqlcgen.CreateOrgPermissionsParams) (int, error) {

	inserter := BulkCreateOrgPermissionInserter{
		OrgPermissions:      input,
		OrganizationService: s,
	}

	affectedRows, err := dbutils.InsertWithTransaction(ctx, s.Conn, inserter)
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal,
			"couldn't create the organization permissions",
			err,
		)
	}

	return int(affectedRows), nil
}

func (s *OrganizationService) CreateOrgPermission(ctx context.Context, input *sqlcgen.CreateOrgPermissionParams) (CreateOrgPermissionOutput, error) {

	insertedID, err := s.DB.CreateOrgPermission(ctx, *input)
	if err != nil {
		return CreateOrgPermissionOutput{}, fmt.Errorf("couldn't create the organization permission: %s", err)
	}

	output := CreateOrgPermissionOutput{
		ID: int(insertedID),
		CreateOrgPermissionParams: CreateOrgPermissionParams{
			OrganizationID: int(input.OrganizationID),
			ResourceTypeID: int(input.ResourceTypeID),
			PermissionCode: input.PermissionCode,
		},
	}

	return output, nil
}

func (s *OrganizationService) GetOrgPermission(ctx context.Context, orgId int, limit int, offset int) ([]CreateOrgPermissionOutput, error) {

	input := sqlcgen.GetOrgPermissionParams{
		OrganizationID: int32(orgId),
		Limit:          int32(limit),
		Offset:         int32(offset),
	}

	orgPermissions, err := s.DB.GetOrgPermission(ctx, input)
	if err != nil {
		return []CreateOrgPermissionOutput{}, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the resource types",
			err,
		)
	}

	if len(orgPermissions) == 0 {
		return []CreateOrgPermissionOutput{}, nil
	}

	var output []CreateOrgPermissionOutput
	for _, orgPermission := range orgPermissions {
		output = append(output, CreateOrgPermissionOutput{
			ID: int(orgPermission.ResourceTypeID),
			CreateOrgPermissionParams: CreateOrgPermissionParams{
				OrganizationID: int(orgPermission.OrganizationID),
				ResourceTypeID: int(orgPermission.ResourceTypeID),
				PermissionCode: orgPermission.PermissionCode,
			},
		})
	}

	return output, nil
}

func (s *OrganizationService) DeleteOrgPermission(ctx context.Context, idType string, id int) error {

	switch strings.ToLower(idType) {
	case "organization":
		if err := s.DB.DeleteOrgPermissionByOrgId(ctx, int32(id)); err != nil {
			return server.NewError(
				server.ErrorInternal,
				"couldn't delete the resource permission by organization_id",
				err,
			)
		}

	case "resource":
		if err := s.DB.DeleteResourceTypeById(ctx, int32(id)); err != nil {
			return server.NewError(
				server.ErrorInternal,
				"couldn't delete the resource permission by id",
				err,
			)
		}
	}

	return nil
}

func (s *OrganizationService) UpsertOrgPermissions(ctx context.Context, orgID int, input []sqlcgen.CreateOrgPermissionsParams) (int, error) {

	existing, err := s.GetOrgPermission(ctx, orgID, 1000, 0)
	if err != nil {
		return 0, err
	}

	existingMap := make(map[string]bool)
	for _, p := range existing {
		key := fmt.Sprintf("%d|%s", p.ResourceTypeID, p.PermissionCode)
		existingMap[key] = true
	}

	var toInsert []sqlcgen.CreateOrgPermissionsParams
	for _, p := range input {
		key := fmt.Sprintf("%d|%s", p.ResourceTypeID, p.PermissionCode)
		if !existingMap[key] {
			toInsert = append(toInsert, p)
		}
	}

	if len(toInsert) == 0 {
		return 0, nil
	}

	return s.CreateOrgPermissionInBatch(ctx, toInsert)
}
