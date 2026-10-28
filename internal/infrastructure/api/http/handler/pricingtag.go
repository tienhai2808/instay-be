package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	pricingTagUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/pricingtag"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type PricingTagHandler struct {
	pricingTagUC pricingTagUC.PricingTagUseCase
}

func NewPricingTagHandler(pricingTagUC pricingTagUC.PricingTagUseCase) *PricingTagHandler {
	return &PricingTagHandler{pricingTagUC}
}

func (h *PricingTagHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.CreatePricingTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.pricingTagUC.Create(ctx, schema, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusCreated, constants.CodeCreatePricingTagSuccess, "Pricing tag created successfully", nil)
}

func (h *PricingTagHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	view := c.Query("view")

	if view == "basic" {
		pricingTagsRes, err := h.pricingTagUC.ListBasic(ctx, schema)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response.OKResponse(c, gin.H{
			"pricing_tags": pricingTagsRes,
		})
		return
	}

	var query request.PricingTagPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	pricingTagsRes, meta, err := h.pricingTagUC.ListFull(ctx, schema, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"pricing_tags": pricingTagsRes,
		"meta":         meta,
	})
}

func (h *PricingTagHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	pricingTagIDStr := c.Param("id")
	pricingTagID, err := strconv.ParseInt(pricingTagIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.UpdatePricingTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.pricingTagUC.Update(ctx, schema, pricingTagID, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdatePricingTagSuccess, "Pricing tag updated successfully", nil)
}

func (h *PricingTagHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	pricingTagIDStr := c.Param("id")
	pricingTagID, err := strconv.ParseInt(pricingTagIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	if err := h.pricingTagUC.Delete(ctx, schema, pricingTagID); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeDeletePricingTagSuccess, "Pricing tag deleted successfully", nil)
}
