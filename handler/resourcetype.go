package handler

import (
	"fmt"
	"strconv"

	srvErr "github.com/bignyap/go-utilities/server"
	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) CreateResurceTypeInBatchHandler(c *gin.Context) {

	input, err := h.ResourceService.CreateResourceTypeJSONValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.ResourceService.CreateResourceTypeInBatch(c, input)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]int{"affected_rows": output})
}

func (h *AdminHandler) CreateResurceTypeHandler(c *gin.Context) {

	input, err := h.ResourceService.CreateResourceTypeFormValidator(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	insertedID, err := h.ResourceService.CreateResourceType(c.Request.Context(), input)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, insertedID)
}

func (h *AdminHandler) ListResourceTypeHandler(c *gin.Context) {

	limit, offset, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	resourceTypes, err := h.ResourceService.ListResourceType(c.Request.Context(), limit, offset)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, resourceTypes)
}

func (h *AdminHandler) DeleteResourceTypeHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	if err := h.ResourceService.DeleteResourceType(c.Request.Context(), int(id)); err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]string{
		"message": fmt.Sprintf("resource type with ID %d deleted successfully", id),
	})
}
