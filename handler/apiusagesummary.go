package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bignyap/go-admin/database/dbutils"
	"github.com/bignyap/go-admin/database/sqlcgen"
	"github.com/bignyap/go-admin/utils/converter"
	"github.com/bignyap/go-admin/utils/formvalidator"
	"github.com/bignyap/go-admin/utils/misc"
)

type CreateApiUsageSummaryParams struct {
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	TotalCalls     int       `json:"total_calls"`
	TotalCost      float64   `json:"total_cost"`
	SubscriptionId int       `json:"subscription_id"`
	ApiEndpointId  int       `json:"api_endpoint_id"`
	OrganizationId int       `json:"organization_id"`
}

type CreateApiUsageSummaryOutput struct {
	ID int `json:"id"`
	CreateApiUsageSummaryParams
}

type CreateApiUsageSummaryInput interface {
	ToCreateApiUsageSummaryParams() CreateApiUsageSummaryParams
}

type LocalApiUsageSummary struct {
	sqlcgen.ApiUsageSummary
}

type LocalCreateApiUsageSummaryParams struct {
	sqlcgen.CreateApiUsageSummaryParams
}

func (apiSummary LocalApiUsageSummary) ToCreateApiUsageSummaryParams() CreateApiUsageSummaryParams {
	return CreateApiUsageSummaryParams{
		StartDate:      misc.FromUnixTime32(apiSummary.UsageStartDate),
		EndDate:        misc.FromUnixTime32(apiSummary.UsageEndDate),
		TotalCalls:     int(apiSummary.TotalCalls),
		TotalCost:      apiSummary.TotalCost,
		SubscriptionId: int(apiSummary.SubscriptionID),
		ApiEndpointId:  int(apiSummary.ApiEndpointID),
		OrganizationId: int(apiSummary.OrganizationID),
	}
}

func (apiSummary LocalCreateApiUsageSummaryParams) ToCreateApiUsageSummaryParams() CreateApiUsageSummaryParams {
	return CreateApiUsageSummaryParams{
		StartDate:      misc.FromUnixTime32(apiSummary.UsageStartDate),
		EndDate:        misc.FromUnixTime32(apiSummary.UsageEndDate),
		TotalCalls:     int(apiSummary.TotalCalls),
		TotalCost:      apiSummary.TotalCost,
		SubscriptionId: int(apiSummary.SubscriptionID),
		ApiEndpointId:  int(apiSummary.ApiEndpointID),
		OrganizationId: int(apiSummary.OrganizationID),
	}
}

func ToCreateApiUsageSummaryOutput(input sqlcgen.ApiUsageSummary) CreateApiUsageSummaryOutput {
	return CreateApiUsageSummaryOutput{
		ID:                          int(input.OrganizationID),
		CreateApiUsageSummaryParams: LocalApiUsageSummary{input}.ToCreateApiUsageSummaryParams(),
	}
}

func CreateApiUsageSummaryFormValidation(r *http.Request) (*sqlcgen.CreateApiUsageSummaryParams, error) {

	err := formvalidator.ParseFormData(r)
	if err != nil {
		return nil, err
	}

	intFields := []string{"total_calls", "subscription_id", "api_endpoint_id", "organization_id"}
	intParsed, err := formvalidator.ParseIntFromForm(r, intFields)
	if err != nil {
		return nil, err
	}

	floatFields := []string{"total_costs"}
	floatParsed, err := formvalidator.ParseFloatFromForm(r, floatFields)
	if err != nil {
		return nil, err
	}

	dateField := []string{"start_date", "end_date"}
	dateParsed, err := formvalidator.ParseUnixTimeFromForm(r, dateField)
	if err != nil {
		return nil, err
	}

	input := sqlcgen.CreateApiUsageSummaryParams{
		UsageStartDate: int32(dateParsed["start_date"]),
		UsageEndDate:   int32(dateParsed["end_date"]),
		TotalCalls:     int32(intParsed["total_calls"]),
		TotalCost:      floatParsed["total_cost"],
		SubscriptionID: int32(intParsed["subscription_id"]),
		OrganizationID: int32(intParsed["organization_id"]),
		ApiEndpointID:  int32(intParsed["api_endpoint_id"]),
	}

	return &input, nil
}

func CreateApiUsageSummaryJSONValidation(r *http.Request) ([]sqlcgen.CreateApiUsageSummariesParams, error) {

	var inputs []CreateApiUsageSummaryParams

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputs)
	if err != nil {
		return nil, err
	}

	var outputs []sqlcgen.CreateApiUsageSummariesParams

	for _, input := range inputs {
		batchInput := sqlcgen.CreateApiUsageSummariesParams{
			UsageStartDate: int32(*converter.TimePtrToUnixInt(&input.StartDate)),
			UsageEndDate:   int32(*converter.TimePtrToUnixInt(&input.EndDate)),
			TotalCalls:     int32(input.TotalCalls),
			TotalCost:      input.TotalCost,
			SubscriptionID: int32(input.SubscriptionId),
			OrganizationID: int32(input.OrganizationId),
			ApiEndpointID:  int32(input.ApiEndpointId),
		}
		outputs = append(outputs, batchInput)
	}

	return outputs, nil
}

type BulkApiSummaryInserter struct {
	ApiUsageSummaries []sqlcgen.CreateApiUsageSummariesParams
	ApiConfig         *ApiConfig
}

func (input BulkApiSummaryInserter) InsertRows(ctx context.Context, tx *sql.Tx) (int64, error) {

	affectedRows, err := input.ApiConfig.DB.CreateApiUsageSummaries(ctx, input.ApiUsageSummaries)
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (apiCfg *ApiConfig) CreateApiUsageInBatchHander(w http.ResponseWriter, r *http.Request) {

	input, err := CreateApiUsageSummaryJSONValidation(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	inserter := BulkApiSummaryInserter{
		ApiUsageSummaries: input,
		ApiConfig:         apiCfg,
	}

	affectedRows, err := dbutils.InsertWithTransaction(r.Context(), apiCfg.Conn, inserter)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the api usage summaries: %s", err))
		return
	}

	respondWithJSON(w, StatusCreated, map[string]int64{"affected_rows": affectedRows})
}

func (apiCfg *ApiConfig) CreateApiUsageHander(w http.ResponseWriter, r *http.Request) {

	input, err := CreateApiUsageSummaryFormValidation(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	organization, err := apiCfg.DB.CreateApiUsageSummary(r.Context(), *input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the api usage summary: %s", err))
		return
	}

	insertedID, err := organization.LastInsertId()
	if err != nil {
		respondWithError(w, StatusInternalServerError, fmt.Sprintf("couldn't retrieve last insert ID: %s", err))
		return
	}

	apiUsageParams := LocalCreateApiUsageSummaryParams{*input}.ToCreateApiUsageSummaryParams()

	output := CreateApiUsageSummaryOutput{
		ID:                          int(insertedID),
		CreateApiUsageSummaryParams: apiUsageParams,
	}

	respondWithJSON(w, StatusCreated, output)
}

func (apiCfg *ApiConfig) GetApiUsageSummaryByOrgIdHandler(w http.ResponseWriter, r *http.Request) {

	id, err := converter.StrToInt(r.PathValue("organization_id"))
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid organization_id format")
		return
	}

	n, page := ExtractPaginationDetail(w, r)

	input := sqlcgen.GetApiUsageSummaryByOrgIdParams{
		OrganizationID: int32(id),
		Limit:          int32(page),
		Offset:         int32(n),
	}

	apiUsageSummaries, err := apiCfg.DB.GetApiUsageSummaryByOrgId(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the api usage summaries: %s", err))
		return
	}

	var output []CreateApiUsageSummaryOutput

	for _, apiUsageSummary := range apiUsageSummaries {

		output = append(output, ToCreateApiUsageSummaryOutput(apiUsageSummary))
	}

	respondWithJSON(w, StatusOK, output)
}

func (apiCfg *ApiConfig) GetApiUsageSummaryBySubIdHandler(w http.ResponseWriter, r *http.Request) {

	id, err := converter.StrToInt(r.PathValue("subscription_id"))
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid subscription_id format")
		return
	}

	n, page := ExtractPaginationDetail(w, r)

	input := sqlcgen.GetApiUsageSummaryBySubIdParams{
		SubscriptionID: int32(id),
		Limit:          int32(page),
		Offset:         int32(n),
	}

	apiUsageSummaries, err := apiCfg.DB.GetApiUsageSummaryBySubId(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the api usage summaries: %s", err))
		return
	}

	var output []CreateApiUsageSummaryOutput

	for _, apiUsageSummary := range apiUsageSummaries {

		output = append(output, ToCreateApiUsageSummaryOutput(apiUsageSummary))
	}

	respondWithJSON(w, StatusOK, output)
}

func (apiCfg *ApiConfig) GetApiUsageSummaryByEndpointIdHandler(w http.ResponseWriter, r *http.Request) {

	id, err := converter.StrToInt(r.PathValue("endpoint_id"))
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid endpoint_id format")
		return
	}

	n, page := ExtractPaginationDetail(w, r)

	input := sqlcgen.GetApiUsageSummaryByEndpointIdParams{
		ApiEndpointID: int32(id),
		Limit:         int32(page),
		Offset:        int32(n),
	}

	apiUsageSummaries, err := apiCfg.DB.GetApiUsageSummaryByEndpointId(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the api usage summaries: %s", err))
		return
	}

	var output []CreateApiUsageSummaryOutput

	for _, apiUsageSummary := range apiUsageSummaries {

		output = append(output, ToCreateApiUsageSummaryOutput(apiUsageSummary))
	}

	respondWithJSON(w, StatusOK, output)
}
