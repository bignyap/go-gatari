package adminHandler

import (
	"fmt"

	converter "github.com/bignyap/go-utilities/converter"
	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) CreateTierPricingInBatchandler(c *gin.Context) {

	input, err := h.PricingService.CreateTierPricingJSONValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	affectedRows, err := h.PricingService.CreateTierPricingInBatch(c.Request.Context(), input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]int{"affected_rows": affectedRows})
}

func (h *AdminHandler) GetTierPricingByTierIdHandler(c *gin.Context) {

	id, err := extractIntPathParam(c, "tier_id")
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	limit, offset, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.PricingService.GetTierPricingByTierId(c.Request.Context(), id, limit, offset)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, output)
}

func extractIntPathParam(c *gin.Context, key string) (int, error) {
	idStr := c.Param(key)
	if idStr == "" {
		return 0, fmt.Errorf("missing %s parameter", key)
	}
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		return 0, fmt.Errorf("invalid %s parameter format", key)
	}
	return id, nil
}

func (h *AdminHandler) DeleteTierPricingHandler(c *gin.Context) {

	orgId := c.Param("tier_id")
	id := c.Param("Id")

	if orgId != "" {

		id32, err := converter.StrToInt(orgId)
		if err != nil {
			h.ResponseWriter.BadRequest(c, "invalid tier_id format")
			return
		}

		err = h.PricingService.DeleteTierPricing(c.Request.Context(), "tier", int(id32))
		if err != nil {
			h.ResponseWriter.Error(c, err)
			return
		}

		h.ResponseWriter.Success(c, map[string]string{
			"message": fmt.Sprintf("subscription with organization_id %d deleted successfully", int32(id32)),
		})
		return
	}

	if id != "" {

		id32, err := converter.StrToInt(id)
		if err != nil {
			h.ResponseWriter.BadRequest(c, "invalid id format")
			return
		}

		err = h.PricingService.DeleteTierPricing(c.Request.Context(), "id", int(id32))
		if err != nil {
			h.ResponseWriter.Error(c, err)
			return
		}

		h.ResponseWriter.Success(c, map[string]string{
			"message": fmt.Sprintf("subscription with id %d deleted successfully", int32(id32)),
		})
		return
	}
	h.ResponseWriter.BadRequest(c, "invalid request")
}
