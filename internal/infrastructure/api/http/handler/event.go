package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	eventUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/event"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	eventUC eventUC.EventUseCase
}

func NewEventHandler(eventUC eventUC.EventUseCase) *EventHandler {
	return &EventHandler{eventUC}
}

func (h *EventHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	id, err := h.eventUC.Create(ctx, schema, userID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusCreated, constants.CodeCreateEventSuccess, "Event created successfully", gin.H{
		"id": strconv.FormatInt(id, 10),
	})
}

func (h *EventHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	eventRes, err := h.eventUC.GetByID(ctx, schema, eventID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"event": eventRes,
	})
}

func (h *EventHandler) GetBySlug(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	bookingRoomID := c.GetInt64(middleware.CtxBookingRoomID)
	if bookingRoomID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	eventSlug := c.Param("slug")

	eventRes, err := h.eventUC.GetBySlug(ctx, schema, bookingRoomID, eventSlug)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"event": eventRes,
	})
}

func (h *EventHandler) ListForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	bookingRoomID := c.GetInt64(middleware.CtxBookingRoomID)
	if bookingRoomID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var query request.EventPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	eventsRes, meta, err := h.eventUC.ListForGuest(ctx, schema, bookingRoomID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"events": eventsRes,
		"meta":   meta,
	})
}

func (h *EventHandler) ListForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var query request.EventPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	eventsRes, meta, err := h.eventUC.ListForAdmin(ctx, schema, userID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"events": eventsRes,
		"meta":   meta,
	})
}

func (h *EventHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	var req request.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err = h.eventUC.Update(ctx, schema, eventID, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdateEventSuccess, "Event updated successfully", nil)
}

func (h *EventHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	if err := h.eventUC.Delete(ctx, schema, eventID, userID); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeDeleteEventSuccess, "Event deleted successfully", nil)
}

func (h *EventHandler) DeleteBatch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.DeleteBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	rowsDeleted, err := h.eventUC.DeleteBatch(ctx, schema, userID, req.IDs)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeDeleteEventsSuccess, "Events deleted successfully", gin.H{
		"count": rowsDeleted,
	})
}
