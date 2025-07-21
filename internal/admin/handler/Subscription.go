package adminHandler

import (
	"strconv"

	converter "github.com/bignyap/go-utilities/converter"
	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) CreateSubscriptionHandler(c *gin.Context) {

	input, err := h.SubscriptionService.CreateSubscriptionValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	outptut, err := h.SubscriptionService.CreateSubscription(c.Request.Context(), input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, outptut)
}

func (h *AdminHandler) CreateSubscriptionInBatchandler(c *gin.Context) {

	inputs, err := h.SubscriptionService.CreateSubscriptionInBatchValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	affectedRows, err := h.SubscriptionService.CreateSubscriptionInBatch(c.Request.Context(), inputs)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]int{"affected_rows": affectedRows})
}

func (h *AdminHandler) DeleteSubscriptionHandler(c *gin.Context) {

	orgId := c.Param("organization_id")
	id := c.Param("id")

	if orgId != "" {

		id32, err := converter.StrToInt(orgId)
		if err != nil {
			h.ResponseWriter.BadRequest(c, "Invalid organization_id format")
			return
		}

		err = h.SubscriptionService.DeleteSubscription(c.Request.Context(), "organization", id32)
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

		err = h.SubscriptionService.DeleteSubscription(c.Request.Context(), "subscription", id32)
		if err != nil {
			h.ResponseWriter.Error(c, err)
			return
		}

		h.ResponseWriter.Success(c, "deleted successfully")
	}
}

func (h *AdminHandler) GetSubscriptionHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, "invalid id format")
		return
	}

	subscription, err := h.SubscriptionService.GetSubscription(c.Request.Context(), int(id))
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, subscription)
}

func (h *AdminHandler) GetSubscriptionByrgIdHandler(c *gin.Context) {

	orgId, err := strconv.Atoi(c.Param("organization_id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, "invalid organization_id format")
		return
	}

	limit, offset, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	subscriptions, err := h.SubscriptionService.GetSubscriptionByOrgId(c.Request.Context(), orgId, limit, offset)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, subscriptions)
}

func (h *AdminHandler) ListSubscriptionHandler(c *gin.Context) {

	limit, offset, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	subscriptions, err := h.SubscriptionService.ListSubscription(c.Request.Context(), limit, offset)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, subscriptions)
}
