package resource

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/bignyap/go-admin/internal/common"
	"github.com/bignyap/go-admin/internal/database/dbutils"
	"github.com/bignyap/go-admin/internal/database/sqlcgen"
	"github.com/bignyap/go-utilities/server"
)

type BulkRegisterEndpointInserter struct {
	Endpoints       []sqlcgen.RegisterApiEndpointsParams
	ResourceService *ResourceService
}

func (input BulkRegisterEndpointInserter) InsertRows(ctx context.Context, tx pgx.Tx) (int64, error) {
	return input.ResourceService.DB.RegisterApiEndpoints(ctx, input.Endpoints)
}

func (s *ResourceService) RegisterApiEndpoint(ctx context.Context, input *RegisterEndpointParams) (RegisterEndpointOutputs, error) {

	var desc pgtype.Text
	if input.Description != nil {
		desc = pgtype.Text{String: *input.Description, Valid: true}
	}

	params := sqlcgen.RegisterApiEndpointParams{
		EndpointName:        input.Name,
		EndpointDescription: desc,
		HttpMethod:          input.HttpMethod,
		PathTemplate:        input.PathTemplate,
		ResourceTypeID:      input.ResourceTypeID,
	}

	err := s.PubSubClient.Publish(ctx, string(common.EndpointCreated), common.EndpointCreatedEvent{
		Path:   input.PathTemplate,
		Method: input.HttpMethod,
		Code:   input.Name,
	})
	if err != nil {
		return RegisterEndpointOutputs{}, server.NewError(
			server.ErrorInternal,
			"couldn't push to the queue",
			err,
		)
	}

	insertedID, err := s.DB.RegisterApiEndpoint(ctx, params)
	if err != nil {
		return RegisterEndpointOutputs{}, server.NewError(
			server.ErrorInternal,
			"couldn't register the API endpoint",
			err,
		)
	}

	output := RegisterEndpointOutputs{
		ID: int(insertedID),
		RegisterEndpointParams: RegisterEndpointParams{
			Name:           input.Name,
			Description:    input.Description,
			HttpMethod:     input.HttpMethod,
			PathTemplate:   input.PathTemplate,
			ResourceTypeID: input.ResourceTypeID,
		},
	}

	return output, nil
}

func (s *ResourceService) RegisterApiEndpointInBatch(ctx context.Context, inputs []RegisterEndpointParams) (int, error) {

	var batch []sqlcgen.RegisterApiEndpointsParams
	for _, in := range inputs {
		desc := pgtype.Text{Valid: false}
		if in.Description != nil {
			desc = pgtype.Text{String: *in.Description, Valid: true}
		}

		dbIn := sqlcgen.RegisterApiEndpointsParams{
			EndpointName:        in.Name,
			EndpointDescription: desc,
			HttpMethod:          in.HttpMethod,
			PathTemplate:        in.PathTemplate,
			ResourceTypeID:      in.ResourceTypeID,
		}

		err := s.PubSubClient.Publish(ctx, string(common.EndpointCreated), common.EndpointCreatedEvent{
			Path:   in.PathTemplate,
			Method: in.HttpMethod,
			Code:   in.Name,
		})
		if err != nil {
			return 0, server.NewError(
				server.ErrorInternal,
				"couldn't push to the queue",
				err,
			)
		}

		batch = append(batch, dbIn)
	}

	inserter := BulkRegisterEndpointInserter{
		Endpoints:       batch,
		ResourceService: s,
	}

	affectedRows, err := dbutils.InsertWithTransaction(ctx, s.Conn, inserter)
	if err != nil {
		return 0, server.NewError(
			server.ErrorInternal,
			"couldn't register endpoints",
			err,
		)
	}

	return int(affectedRows), nil
}

func (s *ResourceService) ListApiEndpoints(ctx context.Context, limit int, offset int) ([]RegisterEndpointOutputs, error) {

	input := sqlcgen.ListApiEndpointParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	apiEndpoints, err := s.DB.ListApiEndpoint(ctx, input)
	if err != nil {
		return []RegisterEndpointOutputs{}, server.NewError(
			server.ErrorInternal,
			"couldn't retrieve endpoints",
			err,
		)
	}

	if len(apiEndpoints) == 0 {
		return []RegisterEndpointOutputs{}, nil
	}

	var output []RegisterEndpointOutputs
	for _, apiEndpoint := range apiEndpoints {
		var desc *string
		if apiEndpoint.EndpointDescription.Valid {
			desc = &apiEndpoint.EndpointDescription.String
		}

		output = append(output, RegisterEndpointOutputs{
			ID: int(apiEndpoint.ApiEndpointID),
			RegisterEndpointParams: RegisterEndpointParams{
				Name:           apiEndpoint.EndpointName,
				Description:    desc,
				HttpMethod:     apiEndpoint.HttpMethod,
				PathTemplate:   apiEndpoint.PathTemplate,
				ResourceTypeID: apiEndpoint.ResourceTypeID,
			},
		})
	}

	return output, nil
}

func (s *ResourceService) DeleteApiEndpointsById(ctx context.Context, id int) error {

	err := s.DB.DeleteApiEndpointById(ctx, int32(id))
	if err != nil {
		return server.NewError(
			server.ErrorInternal,
			"couldn't delete the endpoint",
			err,
		)
	}

	return nil
}
