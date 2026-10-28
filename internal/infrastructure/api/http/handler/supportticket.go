package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	supportTicketUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/supportticket"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type SupportTicketHandler struct {
	supportTicketUC supportTicketUC.SupportTicketUseCase
}

func NewSupportTicketHandler(supportTicketUC supportTicketUC.SupportTicketUseCase) *SupportTicketHandler {
	return &SupportTicketHandler{supportTicketUC}
}

func (h *SupportTicketHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	bookingRoomID := c.GetInt64(middleware.CtxBookingRoomID)
	if bookingRoomID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.CreateSupportTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	id, err := h.supportTicketUC.Create(ctx, schema, bookingRoomID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusCreated, constants.CodeCreateSupportTicketSuccess, "Support ticket created successfully", gin.H{
		"id": strconv.FormatInt(id, 10),
	})
}

func (h *SupportTicketHandler) ListForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	bookingRoomID := c.GetInt64(middleware.CtxBookingRoomID)
	if bookingRoomID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	supTicketsRes, err := h.supportTicketUC.ListForGuest(ctx, schema, bookingRoomID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"support_tickets": supTicketsRes,
	})
}

func (h *SupportTicketHandler) ListForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	departmentIDAny, exists := c.Get(middleware.CtxDepartmentID)
	if !exists {
		_ = c.Error(exceptions.ErrForbidden)
		return
	}

	deptID, ok := departmentIDAny.(*int64)
	if !ok {
		_ = c.Error(exceptions.ErrForbidden)
		return
	}

	var query request.SupportTicketPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	supTicketsRes, meta, err := h.supportTicketUC.ListForAdmin(ctx, schema, userID, deptID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"support_tickets": supTicketsRes,
		"meta":            meta,
	})
}

func (h *SupportTicketHandler) UpdateStatusForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	supTickIDStr := c.Param("id")
	supTickID, err := strconv.ParseInt(supTickIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	bookingRoomID := c.GetInt64(middleware.CtxBookingRoomID)
	if bookingRoomID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.UpdateSupportTicketStatusForGuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.supportTicketUC.UpdateStatusForGuest(ctx, schema, supTickID, bookingRoomID, req); err != nil {
		_ = c.Error(err)
		return
	}

}

func (h *SupportTicketHandler) UpdateStatusForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	supTickIDStr := c.Param("id")
	supTickID, err := strconv.ParseInt(supTickIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.UpdateSupportTicketStatusForAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.supportTicketUC.UpdateStatusForAdmin(ctx, schema, supTickID, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdateSupportTicketSuccess, "Support ticket updated successfully", nil)
}
