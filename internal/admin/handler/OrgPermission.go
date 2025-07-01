package adminHandler

import (
	"fmt"

	converter "github.com/bignyap/go-utilities/converter"
	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) CreateOrgPermissionInBatchHandler(c *gin.Context) {

	input, err := h.OrganizationService.CreateOrgPermissionJSONValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.OrganizationService.CreateOrgPermissionInBatch(c.Request.Context(), input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]int{"affected_rows": output})
}

func (h *AdminHandler) CreateOrgPermissionHandler(c *gin.Context) {

	input, err := h.OrganizationService.CreateOrgPermissionFormValidator(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	insertedID, err := h.OrganizationService.CreateOrgPermission(c.Request.Context(), input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, insertedID)
}

func (h *AdminHandler) GetOrgPermissionHandler(c *gin.Context) {

	id, err := converter.StrToInt(c.Param("organization_id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, "invalid id format")
		return
	}

	limit, offset, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	orgPermissions, err := h.OrganizationService.GetOrgPermission(c.Request.Context(), id, limit, offset)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, orgPermissions)
}

func (h *AdminHandler) DeleteOrgPermissionHandler(c *gin.Context) {

	orgId := c.Param("organization_id")
	id := c.Param("id")

	if orgId != "" {
		id32, err := converter.StrToInt(orgId)
		if err != nil {
			h.ResponseWriter.BadRequest(c, "invalid organization_id format")
			return
		}
		if err := h.OrganizationService.DeleteOrgPermission(c.Request.Context(), "organization ", int(id32)); err != nil {
			h.ResponseWriter.Error(c, err)
			return
		}
		h.ResponseWriter.Success(c,
			map[string]string{
				"message": fmt.Sprintf("resource permission with organization_id %d deleted successfully", int32(id32)),
			},
		)
		return
	}

	if id != "" {
		id32, err := converter.StrToInt(id)
		if err != nil {
			h.ResponseWriter.BadRequest(c, "invalid id format")
			return
		}
		if err := h.OrganizationService.DeleteOrgPermission(c.Request.Context(), "resource ", int(id32)); err != nil {
			h.ResponseWriter.Error(c, err)
			return
		}
		h.ResponseWriter.Success(c,
			map[string]string{
				"message": fmt.Sprintf("resource permission with id %d deleted successfully", int32(id32)),
			},
		)
		return
	}

	h.ResponseWriter.BadRequest(c, "invalid request")
}

func (h *AdminHandler) UpdateOrgPermissionInBatchHandler(c *gin.Context) {

	input, err := h.OrganizationService.UpdateOrgPermissionJSONValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.OrganizationService.UpsertOrgPermissions(c.Request.Context(), int(input[1].OrganizationID), input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]int{"affected_rows": output})
}
