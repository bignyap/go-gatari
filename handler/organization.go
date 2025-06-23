package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	srvErr "github.com/bignyap/go-utilities/server"
)

func (h *AdminHandler) CreateOrganizationandler(c *gin.Context) {

	input, err := h.OrganizationService.ValidateOrgInput(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.OrganizationService.CreateOrganization(c, input)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, output)
}

func (h *AdminHandler) CreateOrganizationInBatchandler(c *gin.Context) {

	input, err := h.OrganizationService.ValidateOrgBatchInput(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	output, err := h.OrganizationService.CreateOrganizationInBatch(c, input)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, output)
}

// func toText(val *string) pgtype.Text {
// 	if val == nil {
// 		return pgtype.Text{Valid: false}
// 	}
// 	return pgtype.Text{String: *val, Valid: true}
// }

// func toBool(val *bool) pgtype.Bool {
// 	if val == nil {
// 		return pgtype.Bool{Valid: false}
// 	}
// 	return pgtype.Bool{Bool: *val, Valid: true}
// }

func (h *AdminHandler) ListOrganizationsHandler(c *gin.Context) {

	limit, offset, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	organizations, err := h.OrganizationService.ListOrganizations(c.Request.Context(), limit, offset)
	if err != nil {
		srvErr.ToApiError(c, err)
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
		srvErr.ToApiError(c, err)
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
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, "deleted successfully")
}
