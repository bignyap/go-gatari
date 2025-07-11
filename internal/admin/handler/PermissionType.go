package adminHandler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) CreatePermissionTypeInBatchHandler(c *gin.Context) {

	input, err := h.ResourceService.CreatePermissionTypeJSONValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.ResourceService.CreatePermissionTypeInBatch(c, input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]int{"affected_rows": output})
}

func (h *AdminHandler) CreatePermissionTypeHandler(c *gin.Context) {

	input, err := h.ResourceService.CreatePermissionTypeFormValidator(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	insertedID, err := h.ResourceService.CreatePermissionType(c.Request.Context(), input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, insertedID)
}

func (h *AdminHandler) ListPermissionTypeHandler(c *gin.Context) {

	limit, offset, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	PermissionTypes, err := h.ResourceService.ListPermissionType(c.Request.Context(), limit, offset)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, PermissionTypes)
}

func (h *AdminHandler) DeletePermissionTypeHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	if err := h.ResourceService.DeletePermissionType(c.Request.Context(), int(id)); err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]string{
		"message": fmt.Sprintf("Permission type with ID %d deleted successfully", id),
	})
}
