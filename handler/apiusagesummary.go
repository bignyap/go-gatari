package handler

import (
	usage "github.com/bignyap/go-admin/internal/service/Usage"
	"github.com/bignyap/go-admin/utils/converter"
	srvErr "github.com/bignyap/go-utilities/server"
	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) CreateApiUsageInBatchHandler(c *gin.Context) {

	input, err := usage.CreateApiUsageSummaryJSONValidation(c)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	output, err := h.UsageService.CreateApiUsageInBatch(c.Request.Context(), input)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Created(c, map[string]int64{"affected_rows": output})
}

// func (h *AdminHandler) CreateApiUsageHandler(c *gin.Context) {

// 	input, err := usageSrv.CreateApiUsageSummaryFormValidation(c)
// 	if err != nil {
// 		h.ResponseWriter.Error(c, err)
// 		return
// 	}

// 	output, err := usageSrv.CreateApiUsage(c.Request.Context(), input)
// 	if err != nil {
// 		srvErr.ToApiError(err)
// 		return
// 	}

// 	h.ResponseWriter.Created(c, output)
// }

func (h *AdminHandler) GetApiUsageSummaryByOrgIdHandler(c *gin.Context) {

	id, err := converter.StrToInt(c.Param("organization_id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, "Invalid organization_id format")
		return
	}

	n, page, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.UsageService.GetApiUsageSummaryByOrgId(c.Request.Context(), id, n, page)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, output)
}

func (h *AdminHandler) GetApiUsageSummaryBySubIdHandler(c *gin.Context) {

	id, err := converter.StrToInt(c.Param("subscription_id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, "Invalid subscription_id format")
		return
	}

	n, page, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.UsageService.GetApiUsageSummaryBySubId(c.Request.Context(), id, n, page)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, output)
}

func (h *AdminHandler) GetApiUsageSummaryByEndpointIdHandler(c *gin.Context) {

	id, err := converter.StrToInt(c.Param("endpoint_id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, "Invalid endpoint_id format")
		return
	}

	n, page, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.UsageService.GetApiUsageSummaryByEndpointId(c.Request.Context(), id, n, page)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, output)
}
