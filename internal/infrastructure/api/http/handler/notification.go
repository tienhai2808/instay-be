package handler

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	notificationUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/notification"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationUC notificationUC.NotificationUseCase
}

func NewNotificationHandler(notificationUC notificationUC.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{notificationUC}
}

func (h *NotificationHandler) CountUnreadForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	deptIDAny, _ := c.Get(middleware.CtxDepartmentID)
	deptID, _ := deptIDAny.(*int64)
	if deptID == nil {
		_ = c.Error(exceptions.ErrForbidden)
		return
	}

	count, err := h.notificationUC.CountUnreadForAdmin(ctx, schema, userID, *deptID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"count": count,
	})
}

func (h *NotificationHandler) CountUnreadForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	bookingRoomID := c.GetInt64(middleware.CtxBookingRoomID)
	if bookingRoomID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	count, err := h.notificationUC.CountUnreadForGuest(ctx, schema, bookingRoomID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"count": count,
	})
}

func (h *NotificationHandler) ListForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	deptIDAny, _ := c.Get(middleware.CtxDepartmentID)
	deptID, _ := deptIDAny.(*int64)
	if deptID == nil {
		_ = c.Error(exceptions.ErrForbidden)
		return
	}

	var query request.NotificationPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	notisRes, meta, err := h.notificationUC.ListForAdmin(ctx, schema, userID, *deptID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"notifications": notisRes,
		"meta":          meta,
	})
}

func (h *NotificationHandler) ListForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	bookingRoomID := c.GetInt64(middleware.CtxBookingRoomID)
	if bookingRoomID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var query request.NotificationPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	notisRes, meta, err := h.notificationUC.ListForGuest(ctx, schema, bookingRoomID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"notifications": notisRes,
		"meta":          meta,
	})
}
