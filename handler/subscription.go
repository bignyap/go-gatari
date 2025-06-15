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
	"github.com/jinzhu/copier"
)

type CreateSubscriptionParams struct {
	Name               string     `json:"name"`
	Type               string     `json:"type"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	StartDate          time.Time  `json:"start_date"`
	APILimit           *int       `json:"api_limit"`
	ExpiryDate         *time.Time `json:"expiry_date"`
	Description        *string    `json:"description"`
	Status             *bool      `json:"status"`
	OrganizationID     int        `json:"organization_id"`
	SubscriptionTierID int        `json:"subscription_tier_id"`
}

type CreateSubscriptionOutput struct {
	ID int `json:"id"`
	CreateSubscriptionParams
}

type ListSubscriptionOutput struct {
	ID       int    `json:"id"`
	TierName string `json:"tier_name"`
	CreateSubscriptionParams
}

type ListSubscriptionOutputWithCount struct {
	TotalItems int                      `json:"total_items"`
	Data       []ListSubscriptionOutput `json:"data"`
}

type CreateSubscriptionInputs interface {
	ToCreateSubscriptionParams() CreateSubscriptionParams
}

type LocalSubscription struct {
	sqlcgen.Subscription
}

type LocalCreateSubscriptionParams struct {
	sqlcgen.CreateSubscriptionParams
}

func (input LocalSubscription) ToCreateSubscriptionParams() CreateSubscriptionParams {
	return CreateSubscriptionParams{
		Name:               input.SubscriptionName,
		Type:               input.SubscriptionType,
		CreatedAt:          misc.FromUnixTime32(input.SubscriptionCreatedDate),
		UpdatedAt:          misc.FromUnixTime32(input.SubscriptionUpdatedDate),
		StartDate:          misc.FromUnixTime32(input.SubscriptionStartDate),
		APILimit:           converter.NullInt32ToInt(&input.SubscriptionApiLimit),
		ExpiryDate:         converter.NullInt32ToTime(&input.SubscriptionExpiryDate),
		Description:        &input.SubscriptionDescription.String,
		Status:             converter.NullBoolToBool(&input.SubscriptionStatus),
		OrganizationID:     int(input.OrganizationID),
		SubscriptionTierID: int(input.SubscriptionTierID),
	}
}

func (input LocalCreateSubscriptionParams) ToCreateSubscriptionParams() CreateSubscriptionParams {
	return CreateSubscriptionParams{
		Name:               input.SubscriptionName,
		Type:               input.SubscriptionType,
		CreatedAt:          misc.FromUnixTime32(input.SubscriptionCreatedDate),
		UpdatedAt:          misc.FromUnixTime32(input.SubscriptionUpdatedDate),
		StartDate:          misc.FromUnixTime32(input.SubscriptionStartDate),
		APILimit:           converter.NullInt32ToInt(&input.SubscriptionApiLimit),
		ExpiryDate:         converter.NullInt32ToTime(&input.SubscriptionExpiryDate),
		Description:        &input.SubscriptionDescription.String,
		Status:             converter.NullBoolToBool(&input.SubscriptionStatus),
		OrganizationID:     int(input.OrganizationID),
		SubscriptionTierID: int(input.SubscriptionTierID),
	}
}

func ToCreateSubscriptionOutput(input sqlcgen.Subscription) CreateSubscriptionOutput {
	return CreateSubscriptionOutput{
		ID:                       int(input.SubscriptionID),
		CreateSubscriptionParams: LocalSubscription{input}.ToCreateSubscriptionParams(),
	}
}

func CreateSubscriptionFormValidation(r *http.Request) (*sqlcgen.CreateSubscriptionParams, error) {

	err := formvalidator.ParseFormData(r)
	if err != nil {
		return nil, err
	}

	strFields := []string{"name", "type"}
	strParsed, err := formvalidator.ParseStringFromForm(r, strFields)
	if err != nil {
		return nil, err
	}

	intFields := []string{"organization_id", "subscription_tier_id"}
	intParsed, err := formvalidator.ParseIntFromForm(r, intFields)
	if err != nil {
		return nil, err
	}

	nullStrField := []string{"description"}
	nullStrParsed, err := formvalidator.ParseNullStringFromForm(r, nullStrField)
	if err != nil {
		return nil, err
	}

	nullBoolFields := []string{"status"}
	nullBoolPared, err := formvalidator.ParseNullBoolFromForm(r, nullBoolFields)
	if err != nil {
		return nil, err
	}

	dateField := []string{"expiry_date"}
	dateParsed, err := formvalidator.ParseNullUnixTime32FromForm(r, dateField)
	if err != nil {
		return nil, err
	}

	nullInt32Field := []string{"api_limit"}
	nullInt32Parsed, err := formvalidator.ParseNullInt32FromForm(r, nullInt32Field)
	if err != nil {
		return nil, err
	}

	startDate, err := converter.StrToUnixTime(r.FormValue("start_date"))
	if err != nil {
		startDate = int(misc.ToUnixTime())
	}

	input := sqlcgen.CreateSubscriptionParams{
		SubscriptionName:        strParsed["name"],
		SubscriptionType:        strParsed["type"],
		SubscriptionCreatedDate: int32(misc.ToUnixTime()),
		SubscriptionUpdatedDate: int32(misc.ToUnixTime()),
		SubscriptionStartDate:   int32(startDate),
		SubscriptionApiLimit:    nullInt32Parsed["api_limit"],
		SubscriptionExpiryDate:  dateParsed["expiry_date"],
		SubscriptionDescription: nullStrParsed["description"],
		SubscriptionStatus:      nullBoolPared["status"],
		OrganizationID:          int32(intParsed["organization_id"]),
		SubscriptionTierID:      int32(intParsed["subscription_tier_id"]),
	}

	return &input, nil
}

func CreateSubscriptionJSONValidator(r *http.Request) ([]sqlcgen.CreateSubscriptionsParams, error) {

	var inputs []CreateSubscriptionParams

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputs)
	if err != nil {
		return nil, err
	}

	var outputs []sqlcgen.CreateSubscriptionsParams

	currentTime := int32(misc.ToUnixTime())

	for _, input := range inputs {
		batchInput := sqlcgen.CreateSubscriptionsParams{
			SubscriptionName:        input.Name,
			SubscriptionType:        input.Type,
			SubscriptionCreatedDate: currentTime,
			SubscriptionUpdatedDate: currentTime,
			SubscriptionStartDate:   currentTime,
			SubscriptionApiLimit:    converter.IntPtrToNullInt32(input.APILimit),
			SubscriptionExpiryDate:  converter.IntPtrToNullInt32(converter.TimePtrToUnixInt(input.ExpiryDate)),
			SubscriptionDescription: converter.StrToNullStr(*input.Description),
			SubscriptionStatus:      converter.BoolPtrToNullBool(input.Status),
			OrganizationID:          int32(input.OrganizationID),
			SubscriptionTierID:      int32(input.SubscriptionTierID),
		}
		outputs = append(outputs, batchInput)
	}

	return outputs, nil
}

type BulkSubscriptionInserter struct {
	Subscriptions []sqlcgen.CreateSubscriptionsParams
	ApiConfig     *ApiConfig
}

func (input BulkSubscriptionInserter) InsertRows(ctx context.Context, tx *sql.Tx) (int64, error) {

	affectedRows, err := input.ApiConfig.DB.CreateSubscriptions(ctx, input.Subscriptions)
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (apiCfg *ApiConfig) CreateSubscriptionInBatchandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateSubscriptionJSONValidator(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	inserter := BulkSubscriptionInserter{
		Subscriptions: input,
		ApiConfig:     apiCfg,
	}

	affectedRows, err := dbutils.InsertWithTransaction(r.Context(), apiCfg.Conn, inserter)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the subscriptions: %s", err))
		return
	}

	respondWithJSON(w, StatusCreated, map[string]int64{"affected_rows": affectedRows})
}

func (apiCfg *ApiConfig) CreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateSubscriptionFormValidation(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	subscription, err := apiCfg.DB.CreateSubscription(r.Context(), *input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the subscription: %s", err))
		return
	}

	insertedID, err := subscription.LastInsertId()
	if err != nil {
		respondWithError(w, StatusInternalServerError, fmt.Sprintf("couldn't retrieve last insert ID: %s", err))
		return
	}

	subscriptionParams := LocalCreateSubscriptionParams{*input}.ToCreateSubscriptionParams()

	output := CreateSubscriptionOutput{
		ID:                       int(insertedID),
		CreateSubscriptionParams: subscriptionParams,
	}

	respondWithJSON(w, StatusCreated, output)
}

func (apiCfg *ApiConfig) DeleteSubscriptionHandler(w http.ResponseWriter, r *http.Request) {

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

	err = apiCfg.DB.DeleteOrganizationById(r.Context(), int32(id64))
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't delete the organization: %s", err))
		return
	}

	respondWithJSON(w, StatusNoContent, map[string]string{
		"message": fmt.Sprintf("organization with ID %d deleted successfully", int32(id64)),
	})
}

func (apiCfg *ApiConfig) GetSubscriptionHandler(w http.ResponseWriter, r *http.Request) {

	id, err := converter.StrToInt(r.PathValue("id"))
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid ID format")
		return
	}

	subscription, err := apiCfg.DB.GetSubscriptionById(r.Context(), int32(id))
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the subscription: %s", err))
		return
	}

	listSubscriptionRow := sqlcgen.ListSubscriptionRow{}
	copier.Copy(&listSubscriptionRow, &subscription)

	output := ToListSubscriptionOutput(listSubscriptionRow)
	respondWithJSON(w, StatusOK, output)
}

func (apiCfg *ApiConfig) GetSubscriptionByrgIdHandler(w http.ResponseWriter, r *http.Request) {

	orgId, err := converter.StrToInt(r.PathValue("organization_id"))
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid organization_id format")
		return
	}

	limit, offset := ExtractPaginationDetail(w, r)
	input := sqlcgen.GetSubscriptionByOrgIdParams{
		OrganizationID: int32(orgId),
		Limit:          int32(limit),
		Offset:         int32(offset),
	}

	subscriptions, err := apiCfg.DB.GetSubscriptionByOrgId(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the subscriptions: %s", err))
		return
	}

	listSubscriptionRows := make([]sqlcgen.ListSubscriptionRow, len(subscriptions))
	for i, sub := range subscriptions {
		listSubscriptionRows[i] = sqlcgen.ListSubscriptionRow(sub)
	}

	output := ToListSubscriptionOutputWithCount(listSubscriptionRows)
	respondWithJSON(w, StatusOK, output)
}

func (apiCfg *ApiConfig) ListSubscriptionHandler(w http.ResponseWriter, r *http.Request) {

	limit, offset := ExtractPaginationDetail(w, r)
	input := sqlcgen.ListSubscriptionParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	subscriptions, err := apiCfg.DB.ListSubscription(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the subscriptions: %s", err))
		return
	}

	output := ToListSubscriptionOutputWithCount(subscriptions)
	respondWithJSON(w, StatusOK, output)
}

func ToListSubscriptionOutput(input sqlcgen.ListSubscriptionRow) ListSubscriptionOutput {
	return ListSubscriptionOutput{
		ID:       int(input.SubscriptionID),
		TierName: input.TierName,
		CreateSubscriptionParams: CreateSubscriptionParams{
			Name:               input.SubscriptionName,
			Type:               input.SubscriptionType,
			CreatedAt:          misc.FromUnixTime32(input.SubscriptionCreatedDate),
			UpdatedAt:          misc.FromUnixTime32(input.SubscriptionUpdatedDate),
			StartDate:          misc.FromUnixTime32(input.SubscriptionStartDate),
			APILimit:           converter.NullInt32ToInt(&input.SubscriptionApiLimit),
			ExpiryDate:         converter.NullInt32ToTime(&input.SubscriptionExpiryDate),
			Description:        &input.SubscriptionDescription.String,
			Status:             converter.NullBoolToBool(&input.SubscriptionStatus),
			OrganizationID:     int(input.OrganizationID),
			SubscriptionTierID: int(input.SubscriptionTierID),
		},
	}
}

func ToListSubscriptionOutputWithCount(inputs []sqlcgen.ListSubscriptionRow) ListSubscriptionOutputWithCount {
	var data []ListSubscriptionOutput
	for _, input := range inputs {
		data = append(data, ToListSubscriptionOutput(input))
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

	return ListSubscriptionOutputWithCount{
		Data:       data,
		TotalItems: totalItems,
	}
}
