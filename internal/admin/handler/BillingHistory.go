package adminHandler

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"

	billing "github.com/bignyap/go-admin/internal/admin/service/Billing"
)

func (h *AdminHandler) CreateBillingHistoryHandler(c *gin.Context) {

	var input billing.CreateBillingHistoryParams

	err := json.NewDecoder(c.Request.Body).Decode(&input)
	if err != nil {
		h.ResponseWriter.BadRequest(c, "Invalid request payload")
		return
	}

	output, err := h.BillingService.CreateBillingHistory(c.Request.Context(), input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, output)

}

func (h *AdminHandler) CreateBillingHistoryInBatchHandler(c *gin.Context) {

	input, err := h.BillingService.CreateBillingHistoryJSONValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.BillingService.CreateBillingHistoryInBatch(c.Request.Context(), input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]int{"affected_rows": output})
}

func (h *AdminHandler) GetBillingHistoryByOrgIdHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("organization_id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, "invalid organization_id")
		return
	}

	n, page, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.BillingService.GetBillingHistoryByOrgId(c.Request.Context(), id, n, page)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, output)
}

func (h *AdminHandler) GetBillingHistoryBySubIdHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("subscription_id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, "invalid subscription_id")
		return
	}

	n, page, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.BillingService.GetBillingHistoryBySubId(c.Request.Context(), id, n, page)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, output)
}

func (h *AdminHandler) GetBillingHistoryByIdHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("billing_id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, "invalid billing_id")
		return
	}

	n, page, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.BillingService.GetBillingHistoryById(c.Request.Context(), id, n, page)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, output)
}
