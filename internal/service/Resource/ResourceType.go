package resource

import (
	"context"
	"fmt"

	"github.com/bignyap/go-admin/database/dbutils"
	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/jackc/pgx/v5"
)

type BulkCreateResourceTypeInserter struct {
	ResourceType    []sqlcgen.CreateResourceTypesParams
	ResourceService *ResourceService
}

func (input BulkCreateResourceTypeInserter) InsertRows(ctx context.Context, tx pgx.Tx) (int64, error) {
	return input.ResourceService.DB.CreateResourceTypes(ctx, input.ResourceType)
}

func (s *ResourceService) CreateResourceTypeInBatch(ctx context.Context, input []sqlcgen.CreateResourceTypesParams) (int, error) {

	inserter := BulkCreateResourceTypeInserter{
		ResourceType:    input,
		ResourceService: s,
	}

	affectedRows, err := dbutils.InsertWithTransaction(ctx, s.Conn, inserter)
	if err != nil {
		return 0, fmt.Errorf("couldn't create the resource types: %s", err)
	}

	return int(affectedRows), nil
}

func (s *ResourceService) CreateResourceType(ctx context.Context, input *sqlcgen.CreateResourceTypeParams) (CreateResourceTypeOutput, error) {

	insertedID, err := s.DB.CreateResourceType(ctx, *input)
	if err != nil {
		return CreateResourceTypeOutput{}, fmt.Errorf("couldn't create the resource type: %s", err)
	}

	description := (*string)(nil)
	if input.ResourceTypeDescription.Valid {
		description = &input.ResourceTypeDescription.String
	}

	output := CreateResourceTypeOutput{
		ID: int(insertedID),
		CreateResourceTypeParams: CreateResourceTypeParams{
			Name:        input.ResourceTypeName,
			Code:        input.ResourceTypeCode,
			Description: description,
		},
	}

	return output, nil
}

func (s *ResourceService) ListResourceType(ctx context.Context, limit int, offset int) ([]CreateResourceTypeOutput, error) {

	input := sqlcgen.ListResourceTypeParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	resourceTypes, err := s.DB.ListResourceType(ctx, input)
	if err != nil {
		return []CreateResourceTypeOutput{}, fmt.Errorf("couldn't retrieve the resource types: %s", err)
	}

	if len(resourceTypes) == 0 {
		return []CreateResourceTypeOutput{}, nil
	}

	var output []CreateResourceTypeOutput
	for _, resourceType := range resourceTypes {
		description := (*string)(nil)
		if resourceType.ResourceTypeDescription.Valid {
			description = &resourceType.ResourceTypeDescription.String
		}
		output = append(output, CreateResourceTypeOutput{
			ID: int(resourceType.ResourceTypeID),
			CreateResourceTypeParams: CreateResourceTypeParams{
				Name:        resourceType.ResourceTypeName,
				Code:        resourceType.ResourceTypeCode,
				Description: description,
			},
		})
	}

	return output, nil
}

func (s *ResourceService) DeleteResourceType(ctx context.Context, id int) error {

	if err := s.DB.DeleteResourceTypeById(ctx, int32(id)); err != nil {
		return fmt.Errorf("couldn't delete the resource type: %s", err)
	}

	return nil
}
