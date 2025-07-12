package adminHandler

import "github.com/gin-gonic/gin"

func (h *AdminHandler) DashboardCountHandler(c *gin.Context) {

	dashboardCounts, err := h.DashboardService.GetDashboardCounts(c.Request.Context())
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, dashboardCounts)
}

func (h *AdminHandler) DashboardUsageHandler(c *gin.Context) {

	filters, err := h.DashboardService.DashboardUsageQueryValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
	}

	dashboardUsages, err := h.DashboardService.GetDashboardUsage(c.Request.Context(), *filters)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, dashboardUsages)
}
