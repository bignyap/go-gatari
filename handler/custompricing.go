package handler

import (
	converter "github.com/bignyap/go-utilities/converter"
	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) CreateCustomPricingInBatchandler(c *gin.Context) {

	input, err := h.PricingService.CreateCustomPricingJSONValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.PricingService.CreateCustomPricingInBatch(c.Request.Context(), input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]int{"affected_rows": output})
}

func (h *AdminHandler) CreateCustomPricingHandler(c *gin.Context) {

	input, err := h.PricingService.CreateCustomPricingFormValidator(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.PricingService.CreateCustomPricing(c.Request.Context(), input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, output)
}

func (h *AdminHandler) DeleteCustomPricingHandler(c *gin.Context) {

	tierId := c.Param("tier_id")
	id := c.Param("id")

	if tierId != "" {

		id32, err := converter.StrToInt(tierId)
		if err != nil {
			h.ResponseWriter.BadRequest(c, "Invalid tier_id format")
			return
		}

		err = h.PricingService.DeleteCustomPricing(c.Request.Context(), "tier", id32)
		if err != nil {
			h.ResponseWriter.Error(c, err)
			return
		}

		h.ResponseWriter.Success(c, "deleted successfully")
	}

	if id != "" {

		id32, err := converter.StrToInt(id)
		if err != nil {
			h.ResponseWriter.BadRequest(c, "Invalid id format")
			return
		}

		err = h.PricingService.DeleteCustomPricing(c.Request.Context(), "id", id32)
		if err != nil {
			h.ResponseWriter.Error(c, err)
			return
		}

		h.ResponseWriter.Success(c, "deleted successfully")
	}

	h.ResponseWriter.BadRequest(c, "Bad request")

}

func (h *AdminHandler) GetCustomPricingHandler(c *gin.Context) {

	id, err := converter.StrToInt(c.Param("subscription_id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, "Invalid ID format")
		return
	}

	limit, offset, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	customPricings, err := h.PricingService.GetCustomPricing(c.Request.Context(), id, limit, offset)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}
	h.ResponseWriter.Success(c, customPricings)
}
