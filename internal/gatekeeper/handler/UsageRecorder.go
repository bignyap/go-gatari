package gateKeeperHandler

import (
	gatekeeping "github.com/bignyap/go-admin/internal/gatekeeper/service/GateKeeping"
	"github.com/gin-gonic/gin"
)

func (h *GateKeeperHandler) UsageRecorderHandler(c *gin.Context) {
	_, output, err := h.UsageRecorderCore(c)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}
	h.ResponseWriter.Success(c, output)
}

func (h *GateKeeperHandler) UsageRecorderCore(c *gin.Context) (*gatekeeping.RecordUsageInput, float64, error) {
	input, err := h.GateKeepingService.RecordUsageValidator(c)
	if err != nil {
		return nil, 0, err
	}
	output, err := h.GateKeepingService.RecordUsage(c.Request.Context(), input)
	if err != nil {
		return nil, 0, err
	}
	return input, output, nil
}
