package handler

import (
	"fmt"

	"github.com/bignyap/go-admin/utils/converter"
	srvErr "github.com/bignyap/go-utilities/server"
	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) CreateOrgTypeInBatchHandler(c *gin.Context) {

	input, err := h.OrganizationService.CreateOrgTypeJSONValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.OrganizationService.CreateOrgTypeInBatch(c.Request.Context(), input)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]int64{"affected_rows": output})
}

func (h *AdminHandler) CreateOrgTypeHandler(c *gin.Context) {

	name, err := h.OrganizationService.CreateOrgTypeFormValidator(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	insertedID, err := h.OrganizationService.CreateOrgType(c.Request.Context(), name)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, insertedID)
}

func (h *AdminHandler) ListOrgTypeHandler(c *gin.Context) {

	limit, offset, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	orgTypes, err := h.OrganizationService.ListOrgType(c.Request.Context(), limit, offset)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, orgTypes)
}

func (h *AdminHandler) DeleteOrgTypeHandler(c *gin.Context) {

	id, err := converter.StrToInt(c.Param("id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, "invalid id")
		return
	}

	if err := h.OrganizationService.DeleteOrgType(c.Request.Context(), int(id)); err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]string{
		"message": fmt.Sprintf("Organization type with ID %d deleted successfully", int32(id)),
	})
}
