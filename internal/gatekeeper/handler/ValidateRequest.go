package gateKeeperHandler

import (
	gatekeeping "github.com/bignyap/go-admin/internal/gatekeeper/service/GateKeeping"
	"github.com/gin-gonic/gin"
)

func (h *GateKeeperHandler) ValidateRequestHandler(c *gin.Context) {
	output, err := h.ValidateRequestCore(c)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}
	h.ResponseWriter.Success(c, output)
}

func (h *GateKeeperHandler) ValidateRequestCore(c *gin.Context) (*gatekeeping.ValidationRequestOutput, error) {
	input, err := h.GateKeepingService.ValidateRequestHeader(c)
	if err != nil {
		return nil, err
	}
	output, err := h.GateKeepingService.ValidateRequest(c.Request.Context(), input)
	if err != nil {
		return nil, err
	}
	return output, nil
}
