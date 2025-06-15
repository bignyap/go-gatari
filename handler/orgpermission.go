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

type CreateOrgPermissionParams struct {
	ResourceTypeID int    `json:"resource_type_id"`
	OrganizationID int    `json:"organization_id"`
	PermissionCode string `json:"permission_code"`
}

type CreateOrgPermissionOutput struct {
	ID int `json:"id"`
	CreateOrgPermissionParams
}

func CreateOrgPermissionFormValidator(r *http.Request) (*sqlcgen.CreateOrgPermissionParams, error) {

	err := formvalidator.ParseFormData(r)
	if err != nil {
		return nil, err
	}

	strFields := []string{"permission_code"}
	strParsed, err := formvalidator.ParseStringFromForm(r, strFields)
	if err != nil {
		return nil, err
	}

	intFields := []string{"resource_type_id", "organization_id"}
	intParsed, err := formvalidator.ParseIntFromForm(r, intFields)
	if err != nil {
		return nil, err
	}

	input := sqlcgen.CreateOrgPermissionParams{
		OrganizationID: int32(intParsed["organization_id"]),
		ResourceTypeID: int32(intParsed["resource_type_id"]),
		PermissionCode: strParsed["permission_code"],
	}

	return &input, nil
}

func CreateOrgPermissionJSONValidation(r *http.Request) ([]sqlcgen.CreateOrgPermissionsParams, error) {

	var inputs []CreateOrgPermissionParams

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputs)
	if err != nil {
		return nil, err
	}

	var outputs []sqlcgen.CreateOrgPermissionsParams

	for _, input := range inputs {
		batchInput := sqlcgen.CreateOrgPermissionsParams{
			OrganizationID: int32(input.OrganizationID),
			ResourceTypeID: int32(input.ResourceTypeID),
			PermissionCode: input.PermissionCode,
		}
		outputs = append(outputs, batchInput)
	}

	return outputs, nil
}

type BulkCreateOrgPermissionInserter struct {
	OrgPermissions []sqlcgen.CreateOrgPermissionsParams
	ApiConfig      *ApiConfig
}

func (input BulkCreateOrgPermissionInserter) InsertRows(ctx context.Context, tx *sql.Tx) (int64, error) {

	affectedRows, err := input.ApiConfig.DB.CreateOrgPermissions(ctx, input.OrgPermissions)
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (apiCfg *ApiConfig) CreateOrgPermissionInBatchHandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateOrgPermissionJSONValidation(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	inserter := BulkCreateOrgPermissionInserter{
		OrgPermissions: input,
		ApiConfig:      apiCfg,
	}

	affectedRows, err := dbutils.InsertWithTransaction(r.Context(), apiCfg.Conn, inserter)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the organization permissions: %s", err))
		return
	}

	respondWithJSON(w, StatusCreated, map[string]int64{"affected_rows": affectedRows})
}

func (apiCfg *ApiConfig) CreateOrgPermissionHandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateOrgPermissionFormValidator(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
	}

	orgPermission, err := apiCfg.DB.CreateOrgPermission(r.Context(), *input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the organization permission: %s", err))
		return
	}

	insertedID, err := orgPermission.LastInsertId()
	if err != nil {
		respondWithError(w, StatusInternalServerError, fmt.Sprintf("couldn't retrieve last insert ID: %s", err))
		return
	}

	output := CreateOrgPermissionOutput{
		ID: int(insertedID),
		CreateOrgPermissionParams: CreateOrgPermissionParams{
			OrganizationID: int(input.OrganizationID),
			ResourceTypeID: int(input.ResourceTypeID),
			PermissionCode: input.PermissionCode,
		},
	}

	respondWithJSON(w, StatusCreated, output)
}

func (apiCfg *ApiConfig) GetOrgPermissionHandler(w http.ResponseWriter, r *http.Request) {

	id, err := converter.StrToInt(r.PathValue("organization_id"))
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid ID format")
		return
	}

	limit, offset := ExtractPaginationDetail(w, r)
	input := sqlcgen.GetOrgPermissionParams{
		OrganizationID: int32(id),
		Limit:          int32(limit),
		Offset:         int32(offset),
	}

	orgPermissions, err := apiCfg.DB.GetOrgPermission(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the resource types: %s", err))
		return
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

	respondWithJSON(w, StatusOK, output)
}

func (apiCfg *ApiConfig) DeleteOrgPermissionHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("organization_id")
	var err error

	if idStr != "" {
		id32, err := converter.StrToInt(idStr)
		if err != nil {
			respondWithError(w, StatusBadRequest, "Invalid organization_id format")
			return
		}

		err = apiCfg.DB.DeleteOrgPermissionByOrgId(r.Context(), int32(id32))
		if err != nil {
			respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't delete the resource permission by organization_id: %s", err))
			return
		}

		respondWithJSON(w, StatusNoContent, map[string]string{
			"message": fmt.Sprintf("resource permission with organization_id %d deleted successfully", int32(id32)),
		})
		return
	}

	idStr = r.PathValue("id")
	if idStr == "" {
		respondWithError(w, StatusBadRequest, "Missing organization_id or id")
		return
	}

	id32, err := converter.StrToInt(idStr)
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid id format")
		return
	}

	err = apiCfg.DB.DeleteResourceTypeById(r.Context(), int32(id32))
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't delete the resource permission by id: %s", err))
		return
	}

	respondWithJSON(w, StatusNoContent, map[string]string{
		"message": fmt.Sprintf("resource permission with id %d deleted successfully", int32(id32)),
	})
}
