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
