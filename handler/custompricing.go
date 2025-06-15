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

type CreateCustomPricingParams struct {
	CustomCostPerCall float64 `json:"custom_cost_per_call"`
	CustomRateLimit   int     `json:"custom_rate_limit"`
	SubscriptionID    int     `json:"subscription_id"`
	TierBasePricingID int     `json:"tier_base_pricing_id"`
}

type CreateCustomPricingOutput struct {
	ID int `json:"id"`
	CreateCustomPricingParams
}

func CreateCustomPricingFormValidator(r *http.Request) (*sqlcgen.CreateCustomPricingParams, error) {

	err := formvalidator.ParseFormData(r)
	if err != nil {
		return nil, err
	}

	intFields := []string{"tier_base_pricing_id", "subscription_id", "custom_rate_limit"}
	intParsed, err := formvalidator.ParseIntFromForm(r, intFields)
	if err != nil {
		return nil, err
	}

	floatFields := []string{"custom_cost_per_call"}
	floatParsed, err := formvalidator.ParseFloatFromForm(r, floatFields)
	if err != nil {
		return nil, err
	}

	input := sqlcgen.CreateCustomPricingParams{
		CustomCostPerCall: floatParsed["custom_cost_per_call"],
		CustomRateLimit:   int32(intParsed["custom_rate_limit"]),
		SubscriptionID:    int32(intParsed["subscription_id"]),
		TierBasePricingID: int32(intParsed["tier_base_pricing_id"]),
	}

	return &input, nil
}

func CreateCustomPricingJSONValidation(r *http.Request) ([]sqlcgen.CreateCustomPricingsParams, error) {

	var inputs []CreateCustomPricingParams

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputs)
	if err != nil {
		return nil, err
	}

	var outputs []sqlcgen.CreateCustomPricingsParams

	for _, input := range inputs {
		batchInput := sqlcgen.CreateCustomPricingsParams{
			CustomCostPerCall: input.CustomCostPerCall,
			CustomRateLimit:   int32(input.CustomRateLimit),
			SubscriptionID:    int32(input.SubscriptionID),
			TierBasePricingID: int32(input.TierBasePricingID),
		}
		outputs = append(outputs, batchInput)
	}

	return outputs, nil
}

type BulkCreateCustomPricingsInserter struct {
	CustomPricings []sqlcgen.CreateCustomPricingsParams
	ApiConfig      *ApiConfig
}

func (input BulkCreateCustomPricingsInserter) InsertRows(ctx context.Context, tx *sql.Tx) (int64, error) {

	affectedRows, err := input.ApiConfig.DB.CreateCustomPricings(ctx, input.CustomPricings)
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (apiCfg *ApiConfig) CreateCustomPricingInBatchandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateCustomPricingJSONValidation(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	inserter := BulkCreateCustomPricingsInserter{
		CustomPricings: input,
		ApiConfig:      apiCfg,
	}

	affectedRows, err := dbutils.InsertWithTransaction(r.Context(), apiCfg.Conn, inserter)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the tier pricings: %s", err))
		return
	}

	respondWithJSON(w, StatusCreated, map[string]int64{"affected_rows": affectedRows})
}

func (apiCfg *ApiConfig) CreateCustomPricingHandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateCustomPricingFormValidator(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	customPrice, err := apiCfg.DB.CreateCustomPricing(r.Context(), *input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the custom pricing: %s", err))
		return
	}

	insertedID, err := customPrice.LastInsertId()
	if err != nil {
		respondWithError(w, StatusInternalServerError, fmt.Sprintf("couldn't retrieve last insert ID: %s", err))
		return
	}

	output := CreateCustomPricingOutput{
		ID: int(insertedID),
		CreateCustomPricingParams: CreateCustomPricingParams{
			CustomCostPerCall: input.CustomCostPerCall,
			CustomRateLimit:   int(input.CustomRateLimit),
			SubscriptionID:    int(input.SubscriptionID),
			TierBasePricingID: int(input.TierBasePricingID),
		},
	}

	respondWithJSON(w, StatusCreated, output)
}

func (apiCfg *ApiConfig) DeleteCustomPricingHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("subscription_id")
	var err error

	if idStr != "" {
		id32, err := converter.StrToInt(idStr)
		if err != nil {
			respondWithError(w, StatusBadRequest, "Invalid subscription_id format")
			return
		}

		err = apiCfg.DB.DeleteCustomPricingBySubscriptionId(r.Context(), int32(id32))
		if err != nil {
			respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't delete the custom pricing by subscription_id: %s", err))
			return
		}

		respondWithJSON(w, StatusNoContent, map[string]string{
			"message": fmt.Sprintf("custom pricing with subscription_id %d deleted successfully", int32(id32)),
		})
		return
	}

	idStr = r.PathValue("id")
	if idStr == "" {
		respondWithError(w, StatusBadRequest, "Missing subscription_id or id")
		return
	}

	id32, err := converter.StrToInt(idStr)
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid id format")
		return
	}

	err = apiCfg.DB.DeleteCustomPricingById(r.Context(), int32(id32))
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't delete the custom pricing by id: %s", err))
		return
	}

	respondWithJSON(w, StatusNoContent, map[string]string{
		"message": fmt.Sprintf("custom pricing with id %d deleted successfully", int32(id32)),
	})
}

func (apiCfg *ApiConfig) GetCustomPricingHandler(w http.ResponseWriter, r *http.Request) {

	id, err := converter.StrToInt(r.PathValue("subscription_id"))
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid ID format")
		return
	}

	limit, offset := ExtractPaginationDetail(w, r)
	input := sqlcgen.GetCustomPricingParams{
		SubscriptionID: int32(id),
		Limit:          int32(limit),
		Offset:         int32(offset),
	}

	customPricings, err := apiCfg.DB.GetCustomPricing(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the custom pricing list: %s", err))
		return
	}

	var output []CreateCustomPricingOutput

	for _, customPricing := range customPricings {

		output = append(output, CreateCustomPricingOutput{
			ID: int(customPricing.TierBasePricingID),
			CreateCustomPricingParams: CreateCustomPricingParams{
				TierBasePricingID: int(customPricing.TierBasePricingID),
				SubscriptionID:    int(customPricing.SubscriptionID),
				CustomCostPerCall: customPricing.CustomCostPerCall,
				CustomRateLimit:   int(customPricing.CustomRateLimit),
			},
		})
	}

	respondWithJSON(w, StatusOK, output)
}
