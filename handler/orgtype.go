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

type CreateOrgTypeInput struct {
	Name string `json:"name"`
}

type CreateOrgTypeOutput struct {
	ID int `json:"id"`
	CreateOrgTypeInput
}

func CreateOrgTypeFormValidator(r *http.Request) (string, error) {

	err := formvalidator.ParseFormData(r)
	if err != nil {
		return "", err
	}

	strFields := []string{"name"}
	strParsed, err := formvalidator.ParseStringFromForm(r, strFields)
	if err != nil {
		return "", err
	}

	return strParsed["name"], nil
}

type CreateOrgTypeParams struct {
	Names []string `json:"name"`
}

func CreateOrgTypeJSONValidation(r *http.Request) (CreateOrgTypeParams, error) {

	var inputs []CreateOrgTypeInput

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputs)
	if err != nil {
		return CreateOrgTypeParams{}, err
	}

	var names []string
	for _, val := range inputs {
		names = append(names, val.Name)
	}

	output := CreateOrgTypeParams{
		Names: names,
	}

	return output, nil
}

type BulkCreateOrgTypeInserter struct {
	OrgTypes  CreateOrgTypeParams
	ApiConfig *ApiConfig
}

func (input BulkCreateOrgTypeInserter) InsertRows(ctx context.Context, tx *sql.Tx) (int64, error) {

	affectedRows, err := input.ApiConfig.DB.CreateOrgTypes(ctx, input.OrgTypes.Names)
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (apiCfg *ApiConfig) CreateOrgTypeInBatchHandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateOrgTypeJSONValidation(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	inserter := BulkCreateOrgTypeInserter{
		OrgTypes:  input,
		ApiConfig: apiCfg,
	}

	affectedRows, err := dbutils.InsertWithTransaction(r.Context(), apiCfg.Conn, inserter)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the organization types: %s", err))
		return
	}

	respondWithJSON(w, StatusCreated, map[string]int64{"affected_rows": affectedRows})
}

func (apiCfg *ApiConfig) CreateOrgTypeHandler(w http.ResponseWriter, r *http.Request) {

	name, err := CreateOrgTypeFormValidator(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
	}

	orgType, err := apiCfg.DB.CreateOrgType(r.Context(), name)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the organization type: %s", err))
		return
	}

	insertedID, err := orgType.LastInsertId()
	if err != nil {
		respondWithError(w, StatusInternalServerError, fmt.Sprintf("couldn't retrieve last insert ID: %s", err))
		return
	}

	output := CreateOrgTypeOutput{
		ID: int(insertedID),
		CreateOrgTypeInput: CreateOrgTypeInput{
			Name: name,
		},
	}

	respondWithJSON(w, StatusCreated, output)
}

func (apiCfg *ApiConfig) ListOrgTypeHandler(w http.ResponseWriter, r *http.Request) {

	limit, offset := ExtractPaginationDetail(w, r)
	input := sqlcgen.ListOrgTypeParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	orgTypes, err := apiCfg.DB.ListOrgType(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the organization types: %s", err))
		return
	}

	var output []CreateOrgTypeOutput

	if len(orgTypes) == 0 {
		respondWithJSON(w, StatusOK, []CreateOrgTypeOutput{})
		return
	}

	for _, orgType := range orgTypes {
		output = append(output, CreateOrgTypeOutput{
			ID: int(orgType.OrganizationTypeID),
			CreateOrgTypeInput: CreateOrgTypeInput{
				Name: orgType.OrganizationTypeName,
			},
		})
	}

	respondWithJSON(w, StatusOK, output)
}

func (apiCfg *ApiConfig) DeleteOrgTypeHandler(w http.ResponseWriter, r *http.Request) {

	id, err := converter.StrToInt(r.PathValue("id"))
	if err != nil {
		respondWithError(w, StatusBadRequest, "ID needs to be a positive integer")
		return
	}

	err = apiCfg.DB.DeleteOrgTypeById(r.Context(), int32(id))
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't delete the organization type: %s", err))
		return
	}

	respondWithJSON(w, StatusNoContent, map[string]string{
		"message": fmt.Sprintf("Organization type with ID %d deleted successfully", int32(id)),
	})
}
