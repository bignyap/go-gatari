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

type CreateBillingHistoryParams struct {
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
	TotalAmountDue float64    `json:"total_amount_due"`
	TotalCalls     int        `json:"total_calls"`
	PaymentStatus  string     `json:"payment_status"`
	PaymentDate    *time.Time `json:"payment_date"`
	CreatedAt      time.Time  `json:"created_at"`
	SubscriptionId int        `json:"subscription_id"`
}

type CreateBillingHistoryOutput struct {
	ID int `json:"id"`
	CreateBillingHistoryParams
}

type CreateBillingHistoryInput interface {
	ToCreateBillingHistoryParams() CreateBillingHistoryParams
}

type LocalCreateBillingHistory struct {
	sqlcgen.BillingHistory
}

type LocalCreateBillingHistoryParams struct {
	sqlcgen.CreateBillingHistoryParams
}

func (billingHistory LocalCreateBillingHistory) ToCreateBillingHistoryParams() CreateBillingHistoryParams {
	return CreateBillingHistoryParams{
		StartDate:      misc.FromUnixTime32(billingHistory.BillingStartDate),
		EndDate:        misc.FromUnixTime32(billingHistory.BillingEndDate),
		PaymentDate:    converter.NullInt32ToTime(&billingHistory.PaymentDate),
		CreatedAt:      misc.FromUnixTime32(billingHistory.BillingCreatedAt),
		TotalCalls:     int(billingHistory.TotalCalls),
		TotalAmountDue: billingHistory.TotalAmountDue,
		PaymentStatus:  billingHistory.PaymentStatus,
		SubscriptionId: int(billingHistory.SubscriptionID),
	}
}

func (billingHistory LocalCreateBillingHistoryParams) ToCreateBillingHistoryParams() CreateBillingHistoryParams {
	return CreateBillingHistoryParams{
		StartDate:      misc.FromUnixTime32(billingHistory.BillingStartDate),
		EndDate:        misc.FromUnixTime32(billingHistory.BillingEndDate),
		PaymentDate:    converter.NullInt32ToTime(&billingHistory.PaymentDate),
		CreatedAt:      misc.FromUnixTime32(billingHistory.BillingCreatedAt),
		TotalCalls:     int(billingHistory.TotalCalls),
		TotalAmountDue: billingHistory.TotalAmountDue,
		PaymentStatus:  billingHistory.PaymentStatus,
		SubscriptionId: int(billingHistory.SubscriptionID),
	}
}

func ToCreateBillingHistoryOutput(input sqlcgen.BillingHistory) CreateBillingHistoryOutput {
	return CreateBillingHistoryOutput{
		ID:                         int(input.BillingID),
		CreateBillingHistoryParams: LocalCreateBillingHistory{input}.ToCreateBillingHistoryParams(),
	}
}

func CreateBillingHistoryFormValidation(r *http.Request) (*sqlcgen.CreateBillingHistoryParams, error) {

	err := formvalidator.ParseFormData(r)
	if err != nil {
		return nil, err
	}

	intFields := []string{"total_calls", "subscription_id"}
	intParsed, err := formvalidator.ParseIntFromForm(r, intFields)
	if err != nil {
		return nil, err
	}

	floatFields := []string{"total_amount_due"}
	floatParsed, err := formvalidator.ParseFloatFromForm(r, floatFields)
	if err != nil {
		return nil, err
	}

	dateFields := []string{"billing_start_date", "billing_end_date"}
	dateParsed, err := formvalidator.ParseUnixTimeFromForm(r, dateFields)
	if err != nil {
		return nil, err
	}

	nullDateField := []string{"payment_date"}
	nulldateParsed, err := formvalidator.ParseNullUnixTime32FromForm(r, nullDateField)
	if err != nil {
		return nil, err
	}

	currentTime := int32(misc.ToUnixTime())

	input := sqlcgen.CreateBillingHistoryParams{
		BillingStartDate: int32(dateParsed["billing_start_date"]),
		BillingEndDate:   int32(dateParsed["billing_end_date"]),
		TotalAmountDue:   floatParsed["total_amount_due"],
		TotalCalls:       int32(intParsed["total_calls"]),
		PaymentStatus:    r.FormValue("payment_status"),
		PaymentDate:      nulldateParsed["payment_date"],
		BillingCreatedAt: currentTime,
		SubscriptionID:   int32(intParsed["subscription_id"]),
	}

	return &input, nil
}

func (apiCfg *ApiConfig) CreateBillingHistoryHandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateBillingHistoryFormValidation(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	billingHistory, err := apiCfg.DB.CreateBillingHistory(r.Context(), *input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the billing history: %s", err))
		return
	}

	insertedID, err := billingHistory.LastInsertId()
	if err != nil {
		respondWithError(w, StatusInternalServerError, fmt.Sprintf("couldn't retrieve last insert ID: %s", err))
		return
	}

	billingHistoryParams := LocalCreateBillingHistoryParams{*input}.ToCreateBillingHistoryParams()

	output := CreateBillingHistoryOutput{
		ID:                         int(insertedID),
		CreateBillingHistoryParams: billingHistoryParams,
	}

	respondWithJSON(w, StatusCreated, output)
}

func CreateBillingHistoryJSONValidation(r *http.Request) ([]sqlcgen.CreateBillingHistoriesParams, error) {

	var inputs []CreateBillingHistoryParams

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputs)
	if err != nil {
		return nil, err
	}

	var outputs []sqlcgen.CreateBillingHistoriesParams

	currentTime := int32(misc.ToUnixTime())

	for _, input := range inputs {
		batchInput := sqlcgen.CreateBillingHistoriesParams{
			BillingStartDate: int32(*converter.TimePtrToUnixInt(&input.StartDate)),
			BillingEndDate:   int32(*converter.TimePtrToUnixInt(&input.EndDate)),
			TotalAmountDue:   input.TotalAmountDue,
			TotalCalls:       int32(input.TotalCalls),
			PaymentStatus:    input.PaymentStatus,
			PaymentDate:      converter.IntPtrToNullInt32(converter.TimePtrToUnixInt(input.PaymentDate)),
			BillingCreatedAt: currentTime,
			SubscriptionID:   int32(input.SubscriptionId),
		}
		outputs = append(outputs, batchInput)
	}

	return outputs, nil
}

type BulkBillingHistoryInserter struct {
	BillingHistories []sqlcgen.CreateBillingHistoriesParams
	ApiConfig        *ApiConfig
}

func (input BulkBillingHistoryInserter) InsertRows(ctx context.Context, tx *sql.Tx) (int64, error) {

	affectedRows, err := input.ApiConfig.DB.CreateBillingHistories(ctx, input.BillingHistories)
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (apiCfg *ApiConfig) CreateBillingHistoryInBatchHandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateBillingHistoryJSONValidation(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	inserter := BulkBillingHistoryInserter{
		BillingHistories: input,
		ApiConfig:        apiCfg,
	}

	affectedRows, err := dbutils.InsertWithTransaction(r.Context(), apiCfg.Conn, inserter)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the billing histories: %s", err))
		return
	}

	respondWithJSON(w, StatusCreated, map[string]int64{"affected_rows": affectedRows})
}

func (apiCfg *ApiConfig) GetBillingHistoryByOrgIdHandler(w http.ResponseWriter, r *http.Request) {

	id, err := converter.StrToInt(r.PathValue("organization_id"))
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid organization_id format")
		return
	}

	n, page := ExtractPaginationDetail(w, r)
	input := sqlcgen.GetBillingHistoryByOrgIdParams{
		OrganizationID: int32(id),
		Limit:          int32(page),
		Offset:         int32(n),
	}

	billingHistories, err := apiCfg.DB.GetBillingHistoryByOrgId(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the billing histories: %s", err))
		return
	}

	respondWithJSON(w, StatusOK, billingHistories)
}

func (apiCfg *ApiConfig) GetBillingHistoryBySubIdHandler(w http.ResponseWriter, r *http.Request) {

	id, err := converter.StrToInt(r.PathValue("subscription_id"))
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid subscription_id format")
		return
	}

	n, page := ExtractPaginationDetail(w, r)
	input := sqlcgen.GetBillingHistoryBySubIdParams{
		SubscriptionID: int32(id),
		Limit:          int32(page),
		Offset:         int32(n),
	}

	billingHistories, err := apiCfg.DB.GetBillingHistoryBySubId(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the billing histories: %s", err))
		return
	}

	respondWithJSON(w, StatusOK, billingHistories)
}

func (apiCfg *ApiConfig) GetBillingHistoryByIdHandler(w http.ResponseWriter, r *http.Request) {

	id, err := converter.StrToInt(r.PathValue("billing_id"))
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid billing_id format")
		return
	}

	n, page := ExtractPaginationDetail(w, r)
	input := sqlcgen.GetBillingHistoryByIdParams{
		BillingID: int32(id),
		Limit:     int32(page),
		Offset:    int32(n),
	}

	billingHistory, err := apiCfg.DB.GetBillingHistoryById(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the billing history: %s", err))
		return
	}

	respondWithJSON(w, StatusOK, billingHistory)
}
