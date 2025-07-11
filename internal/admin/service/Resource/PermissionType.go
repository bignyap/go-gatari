package resource

import (
	"context"

	"github.com/bignyap/go-admin/internal/database/dbutils"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/server"
	"github.com/jackc/pgx/v5"
)

type BulkCreatePermissionTypeInserter struct {
	PermissionType  []sqlcgen.CreatePermissionTypesParams
	ResourceService *ResourceService
}

func (input BulkCreatePermissionTypeInserter) InsertRows(ctx context.Context, tx pgx.Tx) (int64, error) {
	return input.ResourceService.DB.CreatePermissionTypes(ctx, input.PermissionType)
}

func (s *ResourceService) CreatePermissionTypeInBatch(ctx context.Context, input []sqlcgen.CreatePermissionTypesParams) (int, error) {

	inserter := BulkCreatePermissionTypeInserter{
		PermissionType:  input,
		ResourceService: s,
	}

	affectedRows, err := dbutils.InsertWithTransaction(ctx, s.Conn, inserter)
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal,
			"couldn't create the permission types",
			err,
		)
	}

	return int(affectedRows), nil
}

func (s *ResourceService) CreatePermissionType(ctx context.Context, input *sqlcgen.CreatePermissionTypeParams) (CreatePermissionTypeOutput, error) {

	_, err := s.DB.CreatePermissionType(ctx, *input)
	if err != nil {
		return CreatePermissionTypeOutput{}, server.NewError(
			server.ErrorInternal,
			"couldn't create the permission type",
			err,
		)
	}

	description := (*string)(nil)
	if input.PermissionDescription.Valid {
		description = &input.PermissionDescription.String
	}

	output := CreatePermissionTypeOutput{
		CreatePermissionTypeParams: CreatePermissionTypeParams{
			Name:        input.PermissionName,
			Code:        input.PermissionCode,
			Description: description,
		},
	}

	return output, nil
}

func (s *ResourceService) ListPermissionType(ctx context.Context, limit int, offset int) ([]CreatePermissionTypeOutput, error) {

	input := sqlcgen.ListPermissionTypesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	PermissionTypes, err := s.DB.ListPermissionTypes(ctx, input)
	if err != nil {
		return []CreatePermissionTypeOutput{}, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve the permission types",
			err,
		)
	}

	if len(PermissionTypes) == 0 {
		return []CreatePermissionTypeOutput{}, nil
	}

	var output []CreatePermissionTypeOutput
	for _, PermissionType := range PermissionTypes {
		description := (*string)(nil)
		if PermissionType.PermissionDescription.Valid {
			description = &PermissionType.PermissionDescription.String
		}
		output = append(output, CreatePermissionTypeOutput{
			CreatePermissionTypeParams: CreatePermissionTypeParams{
				Name:        PermissionType.PermissionName,
				Code:        PermissionType.PermissionCode,
				Description: description,
			},
		})
	}

	return output, nil
}

func (s *ResourceService) DeletePermissionType(ctx context.Context, id int) error {

	if err := s.DB.DeleteOrgPermissionById(ctx, int32(id)); err != nil {
		return server.NewError(
			server.ErrorInternal,
			"couldn't delete the permission type",
			err,
		)
	}

	return nil
}
