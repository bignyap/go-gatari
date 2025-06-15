package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bignyap/go-admin/database/dbutils"
	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-admin/utils/converter"
	"github.com/bignyap/go-admin/utils/formvalidator"
	"github.com/bignyap/go-admin/utils/misc"
)

type CreateOrganizationParams struct {
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Realm        string    `json:"realm"`
	Country      *string   `json:"country"`
	SupportEmail string    `json:"support_email"`
	Active       *bool     `json:"active"`
	ReportQ      *bool     `json:"report_q"`
	Config       *string   `json:"config"`
	TypeID       int       `json:"type_id"`
}

type CreateOrganizationOutput struct {
	ID int `json:"id"`
	CreateOrganizationParams
}

type ListOrganizationOutput struct {
	ID                   int    `json:"id"`
	OrganizationTypeName string `json:"type"`
	CreateOrganizationParams
}

type ListOrganizationOutputWithCount struct {
	TotalItems int                      `json:"total_items"`
	Data       []ListOrganizationOutput `json:"data"`
}

type CreateOrganizationInputs interface {
	ToCreateOrganizationParams() CreateOrganizationParams
}

type LocalOrganization struct {
	sqlcgen.Organization
}

type LocalCreateCreateOrganizationParams struct {
	sqlcgen.CreateOrganizationParams
}

func (organization LocalOrganization) ToCreateOrganizationParams() CreateOrganizationParams {
	return CreateOrganizationParams{
		Name:         organization.OrganizationName,
		SupportEmail: organization.OrganizationSupportEmail,
		CreatedAt:    misc.FromUnixTime32(organization.OrganizationCreatedAt),
		UpdatedAt:    misc.FromUnixTime32(organization.OrganizationUpdatedAt),
		Realm:        organization.OrganizationRealm,
		Active:       converter.NullBoolToBool(&organization.OrganizationActive),
		ReportQ:      converter.NullBoolToBool(&organization.OrganizationReportQ),
		TypeID:       int(organization.OrganizationTypeID),
		Config:       converter.NullStrToStr(&organization.OrganizationConfig),
	}
}

func (organization LocalCreateCreateOrganizationParams) ToCreateOrganizationParams() CreateOrganizationParams {
	return CreateOrganizationParams{
		Name:         organization.OrganizationName,
		SupportEmail: organization.OrganizationSupportEmail,
		CreatedAt:    misc.FromUnixTime32(organization.OrganizationCreatedAt),
		UpdatedAt:    misc.FromUnixTime32(organization.OrganizationUpdatedAt),
		Realm:        organization.OrganizationRealm,
		Active:       converter.NullBoolToBool(&organization.OrganizationActive),
		ReportQ:      converter.NullBoolToBool(&organization.OrganizationReportQ),
		TypeID:       int(organization.OrganizationTypeID),
		Config:       converter.NullStrToStr(&organization.OrganizationConfig),
	}
}

func ToCreateOrganizationOutput(input sqlcgen.Organization) CreateOrganizationOutput {
	return CreateOrganizationOutput{
		ID:                       int(input.OrganizationID),
		CreateOrganizationParams: LocalOrganization{input}.ToCreateOrganizationParams(),
	}
}

func CreateOrgFormValidation(r *http.Request) (*sqlcgen.CreateOrganizationParams, error) {

	err := formvalidator.ParseFormData(r)
	if err != nil {
		return nil, err
	}

	strFields := []string{"name", "realm", "support_email", "type_id"}
	strParsed, err := formvalidator.ParseStringFromForm(r, strFields)
	if err != nil {
		return nil, err
	}

	intFields := []string{"type_id"}
	intParsed, err := formvalidator.ParseIntFromForm(r, intFields)
	if err != nil {
		return nil, err
	}

	nullStrField := []string{"country", "config"}
	nullStrParsed, err := formvalidator.ParseNullStringFromForm(r, nullStrField)
	if err != nil {
		return nil, err
	}

	boolFields := []string{"active", "reportq"}
	boolParsed, err := formvalidator.ParseNullBoolFromForm(r, boolFields)
	if err != nil {
		return nil, err
	}

	input := sqlcgen.CreateOrganizationParams{
		OrganizationName:         strParsed["name"],
		OrganizationCreatedAt:    int32(misc.ToUnixTime()),
		OrganizationUpdatedAt:    int32(misc.ToUnixTime()),
		OrganizationRealm:        strParsed["realm"],
		OrganizationSupportEmail: strParsed["support_email"],
		OrganizationTypeID:       int32(intParsed["type_id"]),
		OrganizationCountry:      nullStrParsed["country"],
		OrganizationConfig:       nullStrParsed["config"],
		OrganizationActive:       boolParsed["active"],
		OrganizationReportQ:      boolParsed["reqportq"],
	}

	return &input, nil
}

func CreateOrgJSONValidation(r *http.Request) ([]sqlcgen.CreateOrganizationsParams, error) {

	var inputs []CreateOrganizationParams

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputs)
	if err != nil {
		return nil, err
	}

	var outputs []sqlcgen.CreateOrganizationsParams

	currentTime := int32(misc.ToUnixTime())

	for _, input := range inputs {
		batchInput := sqlcgen.CreateOrganizationsParams{
			OrganizationName:         input.Name,
			OrganizationSupportEmail: input.SupportEmail,
			OrganizationRealm:        input.Realm,
			OrganizationActive:       converter.BoolPtrToNullBool(input.Active),
			OrganizationReportQ:      converter.BoolPtrToNullBool(input.ReportQ),
			OrganizationTypeID:       int32(input.TypeID),
			OrganizationConfig:       converter.StrToNullStr(*input.Config),
			OrganizationCreatedAt:    currentTime,
			OrganizationUpdatedAt:    currentTime,
		}
		outputs = append(outputs, batchInput)
	}

	return outputs, nil
}

type BulkOrganizationInserter struct {
	Organizations []sqlcgen.CreateOrganizationsParams
	ApiConfig     *ApiConfig
}

func (input BulkOrganizationInserter) InsertRows(ctx context.Context, tx *sql.Tx) (int64, error) {

	affectedRows, err := input.ApiConfig.DB.CreateOrganizations(ctx, input.Organizations)
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (apiCfg *ApiConfig) CreateOrganizationInBatchandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateOrgJSONValidation(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	inserter := BulkOrganizationInserter{
		Organizations: input,
		ApiConfig:     apiCfg,
	}

	affectedRows, err := dbutils.InsertWithTransaction(r.Context(), apiCfg.Conn, inserter)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the organizations: %s", err))
		return
	}

	respondWithJSON(w, StatusCreated, map[string]int64{"affected_rows": affectedRows})
}

func (apiCfg *ApiConfig) CreateOrganizationandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateOrgFormValidation(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	organization, err := apiCfg.DB.CreateOrganization(r.Context(), *input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the organization: %s", err))
		return
	}

	insertedID, err := organization.LastInsertId()
	if err != nil {
		respondWithError(w, StatusInternalServerError, fmt.Sprintf("couldn't retrieve last insert ID: %s", err))
		return
	}

	organizationParams := LocalCreateCreateOrganizationParams{*input}.ToCreateOrganizationParams()

	output := CreateOrganizationOutput{
		ID:                       int(insertedID),
		CreateOrganizationParams: organizationParams,
	}

	respondWithJSON(w, StatusCreated, output)
}

func (apiCfg *ApiConfig) ListOrganizationsHandler(w http.ResponseWriter, r *http.Request) {

	limit, offset := ExtractPaginationDetail(w, r)
	input := sqlcgen.ListOrganizationParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	organizations, err := apiCfg.DB.ListOrganization(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the organizations: %s", err))
		return
	}

	output := ToListOrganizationOutputWithCount(organizations)
	respondWithJSON(w, StatusOK, output)
}

func ToListOrganizationOutput(input sqlcgen.ListOrganizationRow) ListOrganizationOutput {
	return ListOrganizationOutput{
		ID:                   int(input.OrganizationID),
		OrganizationTypeName: input.OrganizationTypeName,
		CreateOrganizationParams: CreateOrganizationParams{
			Name:         input.OrganizationName,
			SupportEmail: input.OrganizationSupportEmail,
			CreatedAt:    misc.FromUnixTime32(input.OrganizationCreatedAt),
			UpdatedAt:    misc.FromUnixTime32(input.OrganizationUpdatedAt),
			Realm:        input.OrganizationRealm,
			Active:       converter.NullBoolToBool(&input.OrganizationActive),
			ReportQ:      converter.NullBoolToBool(&input.OrganizationReportQ),
			TypeID:       int(input.OrganizationTypeID),
			Config:       converter.NullStrToStr(&input.OrganizationConfig),
			Country:      converter.NullStrToStr(&input.OrganizationCountry),
		},
	}
}

func ToListOrganizationOutputWithCount(inputs []sqlcgen.ListOrganizationRow) ListOrganizationOutputWithCount {
	var data []ListOrganizationOutput
	for _, input := range inputs {
		data = append(data, ToListOrganizationOutput(input))
	}

	totalItems := 0
	if len(inputs) > 0 {
		switch total := inputs[0].TotalItems.(type) {
		case int64:
			totalItems = int(total)
		case int:
			totalItems = total
		default:
			totalItems = 0
		}
	}

	return ListOrganizationOutputWithCount{
		Data:       data,
		TotalItems: totalItems,
	}
}

func (apiCfg *ApiConfig) GetOrganizationByIdHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("Id")
	if idStr == "" {
		respondWithError(w, StatusBadRequest, "ID is required")
		return
	}

	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid ID format")
		return
	}
	input := sqlcgen.ListOrganizationParams{
		Limit:          1,
		Offset:         0,
		OrganizationID: int32(id64),
	}

	organization, err := apiCfg.DB.ListOrganization(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the organization: %s", err))
		return
	}

	if len(organization) == 0 {
		respondWithJSON(w, StatusOK, map[string]interface{}{})
		return
	}

	output := ToListOrganizationOutput(organization[0])
	respondWithJSON(w, StatusOK, output)
}

func (apiCfg *ApiConfig) DeleteOrganizationByIdHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("Id")
	if idStr == "" {
		respondWithError(w, StatusBadRequest, "ID is required")
		return
	}

	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid ID format")
		return
	}

	err = apiCfg.DB.DeleteOrganizationById(r.Context(), int32(id64))
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't delete the organization: %s", err))
		return
	}

	respondWithJSON(w, StatusNoContent, nil)
}
