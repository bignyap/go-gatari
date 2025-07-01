package adminHandler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) CreateOrganizationandler(c *gin.Context) {

	input, err := h.OrganizationService.CreateOrgFormValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.OrganizationService.CreateOrganization(c, input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, output)
}

func (h *AdminHandler) CreateOrganizationInBatchandler(c *gin.Context) {

	input, err := h.OrganizationService.CreateOrgJSONValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.OrganizationService.CreateOrganizationInBatch(c, input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]int{"affected_rows": output})
}

func (h *AdminHandler) ListOrganizationsHandler(c *gin.Context) {

	limit, offset, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	organizations, err := h.OrganizationService.ListOrganizations(c.Request.Context(), limit, offset)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, organizations)
}

func (h *AdminHandler) GetOrganizationByIdHandler(c *gin.Context) {

	id64, err := strconv.ParseInt(c.Param("Id"), 10, 32)
	if err != nil {
		h.ResponseWriter.BadRequest(c, "invalid id format")
		return
	}

	organization, err := h.OrganizationService.GetOrganizationById(c.Request.Context(), int(id64))
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, organization)
}

func (h *AdminHandler) DeleteOrganizationByIdHandler(c *gin.Context) {

	id64, err := strconv.ParseInt(c.Param("Id"), 10, 32)
	if err != nil {
		h.ResponseWriter.BadRequest(c, "invalid id format")
		return
	}

	err = h.OrganizationService.DeleteOrganizationById(c.Request.Context(), int(id64))
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, "deleted successfully")
}

func (h *AdminHandler) UpdateOrganizationandler(c *gin.Context) {

	input, err := h.OrganizationService.UpdateOrgFormValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	err = h.OrganizationService.UpdateOrganization(c, input)
	if err != nil {
		h.ResponseWriter.Error(c, err)
		return
	}

	h.ResponseWriter.Success(c, "organization updated successfully")
}
