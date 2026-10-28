package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	roomTypePriceUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/roomtypeprice"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type RoomTypePriceHandler struct {
	roomTypePriceUC roomTypePriceUC.RoomTypePriceUseCase
}

func NewRoomTypePriceHandler(roomTypePriceUC roomTypePriceUC.RoomTypePriceUseCase) *RoomTypePriceHandler {
	return &RoomTypePriceHandler{roomTypePriceUC}
}

func (h *RoomTypePriceHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var req request.CreateRoomTypePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.roomTypePriceUC.Create(ctx, schema, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusCreated, constants.CodeCreateRoomTypePriceSuccess, "Room type price created successfully", nil)
}

func (h *RoomTypePriceHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	roomTypePriceIDStr := c.Param("id")
	roomTypePriceID, err := strconv.ParseInt(roomTypePriceIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	var req request.UpdateRoomTypePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.roomTypePriceUC.Update(ctx, schema, roomTypePriceID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdateRoomTypePriceSuccess, "Room type price updated successfully", nil)
}

func (h *RoomTypePriceHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	roomTypePriceIDStr := c.Param("id")
	roomTypePriceID, err := strconv.ParseInt(roomTypePriceIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	if err := h.roomTypePriceUC.Delete(ctx, schema, roomTypePriceID); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeDeleteRoomTypePriceSuccess, "Room type price deleted successfully", nil)
}
