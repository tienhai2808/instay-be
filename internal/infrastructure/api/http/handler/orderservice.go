package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	orderServiceUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/orderservice"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type OrderServiceHandler struct {
	orderServiceUC orderServiceUC.OrderServiceUseCase
}

func NewOrderServiceHandler(orderServiceUC orderServiceUC.OrderServiceUseCase) *OrderServiceHandler {
	return &OrderServiceHandler{orderServiceUC}
}

func (h *OrderServiceHandler) CreateForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	bookingRoomID := c.GetInt64(middleware.CtxBookingRoomID)
	if bookingRoomID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.CreateOrderServiceForGuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	id, err := h.orderServiceUC.CreateForGuest(ctx, schema, bookingRoomID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusCreated, constants.CodeCreateOrderServiceSuccess, "Order service created successfully", gin.H{
		"id": strconv.FormatInt(id, 10),
	})
}

func (h *OrderServiceHandler) ListForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	bookingRoomID := c.GetInt64(middleware.CtxBookingRoomID)
	if bookingRoomID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	orderServicesRes, err := h.orderServiceUC.ListForGuest(ctx, schema, bookingRoomID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"order_services": orderServicesRes,
	})
}

func (h *OrderServiceHandler) UpdateStatusForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	orderSvcIDStr := c.Param("id")
	orderSvcID, err := strconv.ParseInt(orderSvcIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	bookingRoomID := c.GetInt64(middleware.CtxBookingRoomID)
	if bookingRoomID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.UpdateOrderServiceStatusForGuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.orderServiceUC.UpdateStatusForGuest(ctx, schema, orderSvcID, bookingRoomID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdateOrderServiceSuccess, "Order service updated successfully", nil)
}

func (h *OrderServiceHandler) CreateForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.CreateOrderServiceForAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	id, err := h.orderServiceUC.CreateForAdmin(ctx, schema, userID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusCreated, constants.CodeCreateOrderServiceSuccess, "Order service created successfully", gin.H{
		"id": strconv.FormatInt(id, 10),
	})
}

func (h *OrderServiceHandler) ListForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

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
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var query request.OrderServicePaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	orderServiceRes, meta, err := h.orderServiceUC.ListForAdmin(ctx, schema, userID, deptID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"order_services": orderServiceRes,
		"meta":           meta,
	})
}

func (h *OrderServiceHandler) GetByIDForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	orderSvcIDStr := c.Param("id")
	orderSvcID, err := strconv.ParseInt(orderSvcIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	orderSvcRes, err := h.orderServiceUC.GetByIDForAdmin(ctx, schema, orderSvcID, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"order_service": orderSvcRes,
	})
}

func (h *OrderServiceHandler) UpdateStatusForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	orderSvcIDStr := c.Param("id")
	orderSvcID, err := strconv.ParseInt(orderSvcIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.UpdateOrderServiceStatusForAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.orderServiceUC.UpdateStatusForAdmin(ctx, schema, orderSvcID, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdateOrderServiceSuccess, "Order service updated successfully", nil)
}
