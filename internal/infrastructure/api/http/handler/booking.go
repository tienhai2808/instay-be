package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	bookingUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/booking"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingUC bookingUC.BookingUseCase
}

func NewBookingHandler(bookingUC bookingUC.BookingUseCase) *BookingHandler {
	return &BookingHandler{bookingUC}
}

func (h *BookingHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	id, err := h.bookingUC.Create(ctx, schema, userID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusCreated, constants.CodeCreateBookingSuccess, "Create booking successfully", gin.H{
		"id": strconv.FormatInt(id, 10),
	})
}

func (h *BookingHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.UpdateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.bookingUC.Update(ctx, schema, bookingID, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdateBookingSuccess, "Update booking successfully", nil)
}

func (h *BookingHandler) UpdateStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.UpdateBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.bookingUC.UpdateStatus(ctx, schema, bookingID, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdateBookingSuccess, "Update booking successfully", nil)
}

func (h *BookingHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	if err = h.bookingUC.Delete(ctx, schema, bookingID, userID); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeDeleteBookingSuccess, "Delete booking successfully", nil)
}

func (h *BookingHandler) StatusSummary(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	statusSummaryRes, err := h.bookingUC.StatusSummary(ctx, schema, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, statusSummaryRes)
}

func (h *BookingHandler) Timeline(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var query request.BookingTimelineQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	timelineRes, err := h.bookingUC.Timeline(ctx, schema, userID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, timelineRes)
}

func (h *BookingHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	bookingIDStr := c.Param("id")
	bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	bookingRes, err := h.bookingUC.GetByID(ctx, schema, bookingID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"booking": bookingRes,
	})
}

func (h *BookingHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var query request.BookingPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	bookingRes, meta, err := h.bookingUC.List(ctx, schema, userID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"bookings": bookingRes,
		"meta":     meta,
	})
}
