package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bignyap/go-admin/database/dbutils"
	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-admin/utils/converter"
	"github.com/bignyap/go-admin/utils/formvalidator"
)

type RegisterEndpointParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type RegisterEndpointOutputs struct {
	ID int `json:"id"`
	RegisterEndpointParams
}

func RegisterEndpointFormValidator(r *http.Request) (*sqlcgen.RegisterApiEndpointParams, error) {

	err := formvalidator.ParseFormData(r)
	if err != nil {
		return nil, err
	}

	strFields := []string{"name"}
	strParsed, err := formvalidator.ParseStringFromForm(r, strFields)
	if err != nil {
		return nil, err
	}

	nullStrFields := []string{"description"}
	nullStrParsed, err := formvalidator.ParseNullStringFromForm(r, nullStrFields)
	if err != nil {
		return nil, err
	}

	input := sqlcgen.RegisterApiEndpointParams{
		EndpointName:        strParsed["name"],
		EndpointDescription: nullStrParsed["description"],
	}

	return &input, nil
}

func RegisterEndpointJSONValidation(r *http.Request) ([]sqlcgen.RegisterApiEndpointsParams, error) {

	var inputs []RegisterEndpointParams

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputs)
	if err != nil {
		return nil, err
	}

	var outputs []sqlcgen.RegisterApiEndpointsParams

	for _, input := range inputs {
		batchInput := sqlcgen.RegisterApiEndpointsParams{
			EndpointName:        input.Name,
			EndpointDescription: converter.StrToNullStr(*input.Description),
		}
		outputs = append(outputs, batchInput)
	}

	return outputs, nil
}

type BulkRegisterEndpointInserter struct {
	Endpoints []sqlcgen.RegisterApiEndpointsParams
	ApiConfig *ApiConfig
}

func (input BulkRegisterEndpointInserter) InsertRows(ctx context.Context, tx *sql.Tx) (int64, error) {

	affectedRows, err := input.ApiConfig.DB.RegisterApiEndpoints(ctx, input.Endpoints)
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (apiCfg *ApiConfig) RegisterEndpointInBatchHandler(w http.ResponseWriter, r *http.Request) {

	input, err := RegisterEndpointJSONValidation(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	fmt.Println(input)

	inserter := BulkRegisterEndpointInserter{
		Endpoints: input,
		ApiConfig: apiCfg,
	}

	affectedRows, err := dbutils.InsertWithTransaction(r.Context(), apiCfg.Conn, inserter)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the endpoints: %s", err))
		return
	}

	respondWithJSON(w, StatusCreated, map[string]int64{"affected_rows": affectedRows})
}

func (apiCfg *ApiConfig) RegisterEndpointHandler(w http.ResponseWriter, r *http.Request) {

	input, err := RegisterEndpointFormValidator(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	apiEndpoint, err := apiCfg.DB.RegisterApiEndpoint(r.Context(), *input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't register the api endpoint: %s", err))
		return
	}

	insertedID, err := apiEndpoint.LastInsertId()
	if err != nil {
		respondWithError(w, StatusInternalServerError, fmt.Sprintf("couldn't retrieve last insert ID: %s", err))
		return
	}

	var description *string
	if input.EndpointDescription.Valid {
		description = &input.EndpointDescription.String
	}

	output := RegisterEndpointOutputs{
		ID: int(insertedID),
		RegisterEndpointParams: RegisterEndpointParams{
			Name:        input.EndpointName,
			Description: description,
		},
	}

	respondWithJSON(w, StatusCreated, output)
}

func (apiCfg *ApiConfig) ListEndpointsHandler(w http.ResponseWriter, r *http.Request) {

	n, page := ExtractPaginationDetail(w, r)
	input := sqlcgen.ListApiEndpointParams{
		Limit:  int32(n),
		Offset: int32(page),
	}

	apiEndpoints, err := apiCfg.DB.ListApiEndpoint(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the organization types: %s", err))
		return
	}

	var output []RegisterEndpointOutputs

	for _, apiEndpoint := range apiEndpoints {

		var description *string
		if apiEndpoint.EndpointDescription.Valid {
			description = &apiEndpoint.EndpointDescription.String
		}

		output = append(output, RegisterEndpointOutputs{
			ID: int(apiEndpoint.ApiEndpointID),
			RegisterEndpointParams: RegisterEndpointParams{
				Name:        apiEndpoint.EndpointName,
				Description: description,
			},
		})
	}

	respondWithJSON(w, StatusOK, output)
}

func (apiCfg *ApiConfig) DeleteEndpointsByIdHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, StatusBadRequest, "ID is required")
		return
	}

	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid ID format")
		return
	}

	err = apiCfg.DB.DeleteApiEndpointById(r.Context(), int32(id64))
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't delete the endpoint: %s", err))
		return
	}

	respondWithJSON(w, StatusNoContent, map[string]string{
		"message": fmt.Sprintf("Api endpoint with ID %d deleted successfully", int32(id64)),
	})
}
