package adminHandler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) CreateSubscriptionTierInBatchHandler(c *gin.Context) {

	input, err := h.SubscriptionService.CreateSubscriptionTierJSONValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	affectedRows, err := h.SubscriptionService.CreateSubscriptionTierInBatch(c.Request.Context(), input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]int{"affected_rows": affectedRows})
}

func (h *AdminHandler) CreateSubscriptionTierHandler(c *gin.Context) {

	input, err := h.SubscriptionService.CreateSubscriptionTierValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.SubscriptionService.CreateSubscriptionTier(c.Request.Context(), *input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, output)
}

func (h *AdminHandler) ListSubscriptionTiersHandler(c *gin.Context) {

	limit, offset, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}
	archived := false
	if v := c.Query("include_archived"); v != "" {
		archived, _ = strconv.ParseBool(v)
	}

	subTiers, err := h.SubscriptionService.ListSubscriptionTiers(c.Request.Context(), limit, offset, archived)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, subTiers)
}

func (h *AdminHandler) DeleteSubscriptionTierHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("Id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	err = h.SubscriptionService.DeleteSubscriptionTier(c.Request.Context(), int(id))
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]string{
		"message": fmt.Sprintf("Subscription tier with ID %d deleted successfully", id),
	})
}
