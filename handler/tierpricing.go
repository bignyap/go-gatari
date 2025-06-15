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

type CreateTierPricingParams struct {
	BaseCostPerCall    float64 `json:"base_cost_per_call"`
	BaseRateLimit      *int    `json:"base_rate_limit"`
	ApiEndpointId      int     `json:"api_endpoint_id"`
	SubscriptionTierID int     `json:"subscription_tier_id"`
}

type CreateTierPricingOutput struct {
	ID int `json:"id"`
	CreateTierPricingParams
}

type CreateTierPricingWithTierName struct {
	CreateTierPricingOutput
	EndpointName string `json:"endpoint_name"`
}

type CreateTierPricingOutputWithCount struct {
	TotalItems int                             `json:"total_items"`
	Data       []CreateTierPricingWithTierName `json:"data"`
}

func CreateTierPricingFormValidator(r *http.Request) (*sqlcgen.CreateTierPricingParams, error) {

	err := formvalidator.ParseFormData(r)
	if err != nil {
		return nil, err
	}

	intFields := []string{"subscription_tier_id", "api_endpoint_id"}
	intParsed, err := formvalidator.ParseIntFromForm(r, intFields)
	if err != nil {
		return nil, err
	}

	nullInt32Field := []string{"base_rate_limit"}
	nullInt32Parsed, err := formvalidator.ParseNullInt32FromForm(r, nullInt32Field)
	if err != nil {
		return nil, err
	}

	floatField := []string{"base_cost_per_call"}
	floatParsed, err := formvalidator.ParseFloatFromForm(r, floatField)
	if err != nil {
		return nil, err
	}

	input := sqlcgen.CreateTierPricingParams{
		SubscriptionTierID: int32(intParsed["subscription_tier_id"]),
		ApiEndpointID:      int32(intParsed["api_endpoint_id"]),
		BaseCostPerCall:    floatParsed["base_cost_per_call"],
		BaseRateLimit:      nullInt32Parsed["base_rate_limit"],
	}

	return &input, nil
}

func CreateTierPricingJSONValidation(r *http.Request) ([]sqlcgen.CreateTierPricingsParams, error) {

	var inputs []CreateTierPricingParams

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&inputs)
	if err != nil {
		return nil, err
	}

	var outputs []sqlcgen.CreateTierPricingsParams

	for _, input := range inputs {
		batchInput := sqlcgen.CreateTierPricingsParams{
			SubscriptionTierID: int32(input.SubscriptionTierID),
			ApiEndpointID:      int32(input.ApiEndpointId),
			BaseCostPerCall:    input.BaseCostPerCall,
			BaseRateLimit:      converter.IntPtrToNullInt32(input.BaseRateLimit),
		}
		outputs = append(outputs, batchInput)
	}

	return outputs, nil
}

type BulkCreateTierPricingsInserter struct {
	TierPricings []sqlcgen.CreateTierPricingsParams
	ApiConfig    *ApiConfig
}

func (input BulkCreateTierPricingsInserter) InsertRows(ctx context.Context, tx *sql.Tx) (int64, error) {

	affectedRows, err := input.ApiConfig.DB.CreateTierPricings(ctx, input.TierPricings)
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (apiCfg *ApiConfig) CreateTierPricingInBatchandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateTierPricingJSONValidation(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	inserter := BulkCreateTierPricingsInserter{
		TierPricings: input,
		ApiConfig:    apiCfg,
	}

	affectedRows, err := dbutils.InsertWithTransaction(r.Context(), apiCfg.Conn, inserter)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the tier pricings: %s", err))
		return
	}

	respondWithJSON(w, StatusCreated, map[string]int64{"affected_rows": affectedRows})
}

func (apiCfg *ApiConfig) CreateTierPricingHandler(w http.ResponseWriter, r *http.Request) {

	input, err := CreateTierPricingFormValidator(r)
	if err != nil {
		respondWithError(w, StatusBadRequest, err.Error())
		return
	}

	tierPrice, err := apiCfg.DB.CreateTierPricing(r.Context(), *input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't create the tier pricing: %s", err))
		return
	}

	insertedID, err := tierPrice.LastInsertId()
	if err != nil {
		respondWithError(w, StatusInternalServerError, fmt.Sprintf("couldn't retrieve last insert ID: %s", err))
		return
	}

	output := CreateTierPricingOutput{
		ID: int(insertedID),
		CreateTierPricingParams: CreateTierPricingParams{
			SubscriptionTierID: int(input.SubscriptionTierID),
			ApiEndpointId:      int(input.ApiEndpointID),
			BaseCostPerCall:    input.BaseCostPerCall,
			BaseRateLimit:      converter.NullInt32ToInt(&input.BaseRateLimit),
		},
	}

	respondWithJSON(w, StatusCreated, output)
}

func (apiCfg *ApiConfig) GetTierPricingByTierIdHandler(w http.ResponseWriter, r *http.Request) {

	idStr, err := converter.StrToInt(r.PathValue("tier_id"))
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid ID format")
		return
	}

	limit, offset := ExtractPaginationDetail(w, r)
	input := sqlcgen.GetTierPricingByTierIdParams{
		SubscriptionTierID: int32(idStr),
		Limit:              int32(limit),
		Offset:             int32(offset),
	}

	tierPricings, err := apiCfg.DB.GetTierPricingByTierId(r.Context(), input)
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't retrieve the tier pricing list: %s", err))
		return
	}

	totalItems := 0
	if len(tierPricings) > 0 {
		switch total := tierPricings[0].TotalItems.(type) {
		case int64:
			totalItems = int(total)
		case int:
			totalItems = total
		default:
			totalItems = 0
		}
	}

	var output []CreateTierPricingWithTierName

	for _, tierPricing := range tierPricings {

		output = append(output, CreateTierPricingWithTierName{
			EndpointName: tierPricing.EndpointName,
			CreateTierPricingOutput: CreateTierPricingOutput{
				ID: int(tierPricing.TierBasePricingID),
				CreateTierPricingParams: CreateTierPricingParams{
					SubscriptionTierID: int(tierPricing.SubscriptionTierID),
					ApiEndpointId:      int(tierPricing.ApiEndpointID),
					BaseCostPerCall:    tierPricing.BaseCostPerCall,
					BaseRateLimit:      converter.NullInt32ToInt(&tierPricing.BaseRateLimit),
				},
			},
		})
	}

	respondWithJSON(w, StatusOK, CreateTierPricingOutputWithCount{
		Data:       output,
		TotalItems: totalItems,
	})
}

func (apiCfg *ApiConfig) DeleteTierPricingHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("organization_id")
	var err error

	if idStr != "" {
		id32, err := converter.StrToInt(idStr)
		if err != nil {
			respondWithError(w, StatusBadRequest, "Invalid organization_id format")
			return
		}

		err = apiCfg.DB.DeleteSubscriptionByOrgId(r.Context(), int32(id32))
		if err != nil {
			respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't delete the subscription by organization_id: %s", err))
			return
		}

		respondWithJSON(w, StatusNoContent, map[string]string{
			"message": fmt.Sprintf("subscription with organization_id %d deleted successfully", int32(id32)),
		})
		return
	}

	idStr = r.PathValue("Id")
	if idStr == "" {
		respondWithError(w, StatusBadRequest, "Missing organization_id or id")
		return
	}

	id32, err := converter.StrToInt(idStr)
	if err != nil {
		respondWithError(w, StatusBadRequest, "Invalid id format")
		return
	}

	err = apiCfg.DB.DeleteSubscriptionById(r.Context(), int32(id32))
	if err != nil {
		respondWithError(w, StatusBadRequest, fmt.Sprintf("couldn't delete the subscription by id: %s", err))
		return
	}

	respondWithJSON(w, StatusNoContent, map[string]string{
		"message": fmt.Sprintf("subscription with id %d deleted successfully", int32(id32)),
	})
}
