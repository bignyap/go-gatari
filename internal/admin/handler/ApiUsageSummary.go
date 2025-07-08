package adminHandler

import (
	usage "github.com/bignyap/go-admin/internal/admin/service/Usage"
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
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Created(c, map[string]int64{"affected_rows": output})
}

func (h *AdminHandler) GetApiUsageSummaryHandler(c *gin.Context) {

	var err error
	var output interface{}

	query, err := h.UsageService.UsageSummaryQueryValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	limit, offset, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	input := usage.UsageSummaryFilters{
		Limit:                         int32(limit),
		Offset:                        int32(offset),
		UsageSummaryFilterQueryParams: query,
	}

	if input.GroupBy {
		output, err = h.UsageService.GetUsageSummaryByDay(c.Request.Context(), input)
	} else {
		output, err = h.UsageService.GetUsageSummary(c.Request.Context(), input)
	}
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, output)
}
