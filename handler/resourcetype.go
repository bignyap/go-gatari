package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bignyap/go-admin/database/dbutils"
	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-admin/utils/converter"
	"github.com/bignyap/go-admin/utils/formvalidator"
)

type CreateResourceTypeParams struct {
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description *string `json:"description"`
}

type CreateResourceTypeOutput struct {
	ID int `json:"id"`
	CreateResourceTypeParams
}

func CreateResourceTypeFormValidator(r *http.Request) (*sqlcgen.CreateResourceTypeParams, error) {

	err := formvalidator.ParseFormData(r)
	if err != nil {
		return nil, err
	}

	strFields := []string{"name", "code"}
	strParsed, err := formvalidator.ParseStringFromForm(r, strFields)
	if err != nil {
		return nil, err
	}

	nullStrField := []string{"description"}
	nullStrParsed, err := formvalidator.ParseNullStringFromForm(r, nullStrField)
	if err != nil {
		return nil, err
	}

	input := sqlcgen.CreateResourceTypeParams{
		ResourceTypeName:        strParsed["name"],
		ResourceTypeCode:        strParsed["code"],
		ResourceTypeDescription: nullStrParsed["description"],
	}

	return &input, nil
}

func CreateResourceTypeJSONValidation(r *http.Request) ([]sqlcgen.CreateResourceTypesParams, error) {

	var inputs []CreateResourceTypeParams

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputs)
	if err != nil {
		return nil, err
	}

	var outputs []sqlcgen.CreateResourceTypesParams

	for _, input := range inputs {
		batchInput := sqlcgen.CreateResourceTypesParams{
			ResourceTypeName:        input.Name,
			ResourceTypeCode:        input.Code,
			ResourceTypeDescription: converter.StrToNullStr(*input.Description),
		}
		outputs = append(outputs, batchInput)
	}

	return outputs, nil
}

type BulkCreateResourceTypeInserter struct {
	ResourceType []sqlcgen.CreateResourceTypesParams
	ApiConfig    *ApiConfig
}

func (input BulkCreateResourceTypeInserter) InsertRows(ctx context.Context, tx *sql.Tx) (int64, error) {

	affectedRows, err := input.ApiConfig.DB.CreateResourceTypes(ctx, input.ResourceType)
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (apiCfg *ApiConfig) CreateResurceTypeInBatchHandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateResourceTypeJSONValidation(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	inserter := BulkCreateResourceTypeInserter{
		ResourceType: input,
		ApiConfig:    apiCfg,
	}

	affectedRows, err := dbutils.InsertWithTransaction(r.Context(), apiCfg.Conn, inserter)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the resource types: %s", err))
		return
	}

	respondWithJSON(w, StatusCreated, map[string]int64{"affected_rows": affectedRows})
}

func (apiCfg *ApiConfig) CreateResurceTypeHandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateResourceTypeFormValidator(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
	}

	resoureType, err := apiCfg.DB.CreateResourceType(r.Context(), *input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the resource type: %s", err))
		return
	}

	insertedID, err := resoureType.LastInsertId()
	if err != nil {
		respondWithError(w, StatusInternalServerError, fmt.Sprintf("couldn't retrieve last insert ID: %s", err))
		return
	}

	output := CreateResourceTypeOutput{
		ID: int(insertedID),
		CreateResourceTypeParams: CreateResourceTypeParams{
			Name:        input.ResourceTypeName,
			Code:        input.ResourceTypeCode,
			Description: converter.NullStrToStr(&input.ResourceTypeDescription),
		},
	}

	respondWithJSON(w, StatusCreated, output)
}

func (apiCfg *ApiConfig) ListResourceTypeHandler(w http.ResponseWriter, r *http.Request) {

	limit, offset := ExtractPaginationDetail(w, r)
	input := sqlcgen.ListResourceTypeParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	resourceTypes, err := apiCfg.DB.ListResourceType(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the resource types: %s", err))
		return
	}

	var output []CreateResourceTypeOutput

	if len(resourceTypes) == 0 {
		respondWithJSON(w, StatusOK, []CreateResourceTypeOutput{})
		return
	}

	for _, resourceType := range resourceTypes {
		output = append(output, CreateResourceTypeOutput{
			ID: int(resourceType.ResourceTypeID),
			CreateResourceTypeParams: CreateResourceTypeParams{
				Name:        resourceType.ResourceTypeName,
				Code:        resourceType.ResourceTypeCode,
				Description: converter.NullStrToStr(&resourceType.ResourceTypeDescription),
			},
		})
	}

	respondWithJSON(w, StatusOK, output)
}

func (apiCfg *ApiConfig) DeleteResourceTypeHandler(w http.ResponseWriter, r *http.Request) {

	id, err := converter.StrToInt(r.PathValue("id"))
	if err != nil {
		respondWithError(w, StatusBadRequest, "ID is required")
		return
	}

	err = apiCfg.DB.DeleteResourceTypeById(r.Context(), int32(id))
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't delete the resource type: %s", err))
		return
	}

	respondWithJSON(w, StatusNoContent, map[string]string{
		"message": fmt.Sprintf("resource type with ID %d deleted successfully", int32(id)),
	})
}
