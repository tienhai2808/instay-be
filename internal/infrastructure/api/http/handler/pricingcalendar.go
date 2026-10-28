package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	pricingCalendarUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/pricingcalendar"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"

	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
)

type PricingCalenderHandler struct {
	pricingCalendarUC pricingCalendarUC.PricingCalendarUseCase
}

func NewPricingCalendarHandler(pricingCalendarUC pricingCalendarUC.PricingCalendarUseCase) *PricingCalenderHandler {
	return &PricingCalenderHandler{pricingCalendarUC}
}

func (h *PricingCalenderHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	pricingCalendarsRes, err := h.pricingCalendarUC.List(ctx, schema)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"pricing_calendars": pricingCalendarsRes,
	})
}

func (h *PricingCalenderHandler) UpdateBulk(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var req request.UpdatePricingCalendarsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.pricingCalendarUC.UpdateBulk(ctx, schema, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdatePricingCalendarSuccess, "Pricing calendar updated successfully", nil)
}
