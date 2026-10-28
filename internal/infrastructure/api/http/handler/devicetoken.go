package handler

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	deviceTokenUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/devicetoken"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type DeviceTokenHandler struct {
	deviceTokenUC deviceTokenUC.DeviceTokenUseCase
}

func NewDeviceTokenHandler(deviceTokenUC deviceTokenUC.DeviceTokenUseCase) *DeviceTokenHandler {
	return &DeviceTokenHandler{deviceTokenUC}
}

func (h *DeviceTokenHandler) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var clientType model.DeviceTokenClientType
	clientID := c.GetInt64(middleware.CtxClientID)
	clientTypeRaw := c.GetString(middleware.CtxClientType)
	switch clientTypeRaw {
	case constants.ClientTypeStaff:
		clientType = model.DeviceTokenClientTypeUser
	case constants.ClientTypeGuest:
		clientType = model.DeviceTokenClientTypeBookingRoom
	}

	if clientID == 0 || clientType == "" {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.DeviceTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.deviceTokenUC.Register(ctx, schema, clientID, clientType, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, nil)
}

func (h *DeviceTokenHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var req request.DeviceTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.deviceTokenUC.Delete(ctx, schema, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, nil)
}
