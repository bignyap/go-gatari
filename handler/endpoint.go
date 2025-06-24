package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	srvErr "github.com/bignyap/go-utilities/server"
)

func (h *AdminHandler) RegisterEndpointHandler(c *gin.Context) {

	fmt.Println(c.Request.Body)

	input, err := h.ResourceService.ValidateRegisterInput(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.ResourceService.RegisterApiEndpoint(c.Request.Context(), input)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Created(c, output)
}

func (h *AdminHandler) RegisterEndpointInBatchHandler(c *gin.Context) {

	input, err := h.ResourceService.ValidateRegisterBatchInput(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.ResourceService.RegisterApiEndpointInBatch(c.Request.Context(), input)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Created(c, map[string]int{"affected_rows": output})
}

func (h *AdminHandler) ListEndpointsHandler(c *gin.Context) {

	n, page, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.ResourceService.ListApiEndpoints(c.Request.Context(), n, page)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Created(c, output)
}

func (h *AdminHandler) DeleteEndpointsByIdHandler(c *gin.Context) {

	id64, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		h.ResponseWriter.BadRequest(c, "invalid id format")
		return
	}

	err = h.ResourceService.DeleteApiEndpointsById(c.Request.Context(), int(id64))
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(
		c, map[string]string{"message": fmt.Sprintf("API endpoint with ID %d deleted successfully", int32(id64))},
	)
}
