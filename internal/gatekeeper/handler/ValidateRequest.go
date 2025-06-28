package gateKeeperHandler

import (
	gatekeeping "github.com/bignyap/go-admin/internal/gatekeeper/service/GateKeeping"
	"github.com/gin-gonic/gin"
)

func (h *GateKeeperHandler) ValidateRequestHandler(c *gin.Context) {
	if _, err := h.ValidateRequestCore(c); err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}
	h.ResponseWriter.Success(c, nil)
}

func (h *GateKeeperHandler) ValidateRequestCore(c *gin.Context) (*gatekeeping.ValidateRequestInput, error) {
	input, err := h.GateKeepingService.ValidateRequestHeader(c)
	if err != nil {
		return nil, err
	}
	err = h.GateKeepingService.ValidateRequest(c.Request.Context(), input)
	if err != nil {
		return nil, err
	}
	return input, nil
}
