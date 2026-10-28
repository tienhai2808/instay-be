package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	roomTypeUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/roomtype"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type RoomTypeHandler struct {
	roomTypeUC roomTypeUC.RoomTypeUseCase
}

func NewRoomTypeHandler(roomTypeUC roomTypeUC.RoomTypeUseCase) *RoomTypeHandler {
	return &RoomTypeHandler{roomTypeUC}
}

func (h *RoomTypeHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.CreateRoomTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.roomTypeUC.Create(ctx, schema, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusCreated, constants.CodeCreateRoomTypeSuccess, "Room type created successfully", nil)
}

func (h *RoomTypeHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	view := c.Query("view")
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var query request.RoomTypePaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if view == "basic" {
		roomTypesRes, err := h.roomTypeUC.ListBasic(ctx, schema, userID, query)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response.OKResponse(c, gin.H{
			"room_types": roomTypesRes,
		})
		return
	}

	roomTypesRes, meta, err := h.roomTypeUC.ListFull(ctx, schema, userID, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"room_types": roomTypesRes,
		"meta":       meta,
	})
}

func (h *RoomTypeHandler) ListRoomsByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	roomTypeIDStr := c.Param("id")
	roomTypeID, err := strconv.ParseInt(roomTypeIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	include := c.Query("include")
	includeRoomType := include == "room_type"
	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	roomTypeRes, roomsRes, err := h.roomTypeUC.ListRoomsByID(ctx, schema, userID, roomTypeID, includeRoomType)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res := gin.H{
		"rooms": roomsRes,
	}

	if includeRoomType {
		res["room_types"] = roomTypeRes
	}

	response.OKResponse(c, res)
}

func (h *RoomTypeHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	roomTypeIDStr := c.Param("id")
	roomTypeID, err := strconv.ParseInt(roomTypeIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.UpdateRoomTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.roomTypeUC.Update(ctx, schema, roomTypeID, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdateRoomTypeSuccess, "Room type updated successfully", nil)
}

func (h *RoomTypeHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	roomTypeIDStr := c.Param("id")
	roomTypeID, err := strconv.ParseInt(roomTypeIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	if err := h.roomTypeUC.Delete(ctx, schema, roomTypeID, userID); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeDeleteRoomTypeSuccess, "Room type deleted successfully", nil)
}

func (h *RoomTypeHandler) PriceCalculation(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var req request.CalculateRoomTypeRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	roomTypePriceRes, err := h.roomTypeUC.CalculatePrice(ctx, schema, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"price": roomTypePriceRes,
	})
}
