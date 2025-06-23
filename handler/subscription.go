package handler

import (
	"fmt"
	"strconv"
	"time"

	srvErr "github.com/bignyap/go-utilities/server"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *AdminHandler) CreateSubscriptionHandler(c *gin.Context) {

	input, err := h.SubscriptionService.CreateSubscriptionValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	outptut, err := h.SubscriptionService.CreateSubscription(c.Request.Context(), input)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, outptut)
}

func (h *AdminHandler) CreateSubscriptionInBatchandler(c *gin.Context) {

	inputs, err := h.SubscriptionService.CreateSubscriptionInBatchValidation(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	affectedRows, err := h.SubscriptionService.CreateSubscriptionInBatch(c.Request.Context(), inputs)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]int{"affected_rows": affectedRows})
}

func (h *AdminHandler) DeleteSubscriptionHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, "invalid id format")
		return
	}

	err = h.SubscriptionService.DeleteSubscription(c.Request.Context(), int(id))
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, map[string]string{
		"message": fmt.Sprintf("organization with ID %d deleted successfully", id),
	})
}

func (h *AdminHandler) GetSubscriptionHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, "invalid id format")
		return
	}

	subscription, err := h.SubscriptionService.GetSubscription(c.Request.Context(), int(id))
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, subscription)
}

func (h *AdminHandler) GetSubscriptionByrgIdHandler(c *gin.Context) {

	orgId, err := strconv.Atoi(c.Param("organization_id"))
	if err != nil {
		h.ResponseWriter.BadRequest(c, "invalid organization_id format")
		return
	}

	limit, offset, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	subscriptions, err := h.SubscriptionService.GetSubscriptionByOrgId(c.Request.Context(), orgId, limit, offset)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, subscriptions)
}

func (h *AdminHandler) ListSubscriptionHandler(c *gin.Context) {

	limit, offset, err := ExtractPaginationDetail(c)
	if err != nil {
		h.ResponseWriter.BadRequest(c, err.Error())
		return
	}

	subscriptions, err := h.SubscriptionService.ListSubscription(c.Request.Context(), limit, offset)
	if err != nil {
		srvErr.ToApiError(c, err)
		return
	}

	h.ResponseWriter.Success(c, subscriptions)
}

// -------- pgtype helpers --------

func toPgInt4(ptr *int) pgtype.Int4 {
	if ptr == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*ptr), Valid: true}
}

// func toPgInt4Ptr(value int) pgtype.Int4 {
// 	return pgtype.Int4{Int32: int32(value), Valid: true}
// }

func toPgInt4Ptr(v *int) pgtype.Int4 {
	if v == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*v), Valid: true}
}

func toPgInt4FromTime(t time.Time) pgtype.Int4 {
	return pgtype.Int4{Int32: int32(t.Unix()), Valid: true}
}

func toPgInt4FromTimePtr(ptr *time.Time) pgtype.Int4 {
	if ptr == nil {
		return pgtype.Int4{Valid: false}
	}
	return toPgInt4FromTime(*ptr)
}

func toPgText(ptr *string) pgtype.Text {
	if ptr == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *ptr, Valid: true}
}

func toPgBool(ptr *bool) pgtype.Bool {
	if ptr == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *ptr, Valid: true}
}

func fromPgInt4Ptr(v pgtype.Int4) *int {
	if !v.Valid {
		return nil
	}
	val := int(v.Int32)
	return &val
}

func fromPgInt4TimePtr(v pgtype.Int4) *time.Time {
	if !v.Valid {
		return nil
	}
	t := time.Unix(int64(v.Int32), 0)
	return &t
}

func fromPgText(v pgtype.Text) *string {
	if !v.Valid {
		return nil
	}
	return &v.String
}

func fromPgBool(v pgtype.Bool) *bool {
	if !v.Valid {
		return nil
	}
	return &v.Bool
}
