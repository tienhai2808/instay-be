package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	supportCategoryUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/supportcategory"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type SupportCategoryHandler struct {
	supportCategoryUC supportCategoryUC.SupportCategoryUseCase
}

func NewSupportCategoryHandler(supportCategoryUC supportCategoryUC.SupportCategoryUseCase) *SupportCategoryHandler {
	return &SupportCategoryHandler{supportCategoryUC}
}

func (h *SupportCategoryHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.CreateSupportCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.supportCategoryUC.Create(ctx, schema, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusCreated, constants.CodeCreateSupportCategorySuccess, "Support category created successfully", nil)
}

func (h *SupportCategoryHandler) ListForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	view := c.Query("view")

	if view == "basic" {
		supportCategorysRes, err := h.supportCategoryUC.ListBasicForAdmin(ctx, schema)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response.OKResponse(c, gin.H{
			"support_categories": supportCategorysRes,
		})
		return
	}

	var query request.SupportCategoryPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	supCatsRes, meta, err := h.supportCategoryUC.ListFullForAdmin(ctx, schema, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"support_categories": supCatsRes,
		"meta":               meta,
	})
}

func (h *SupportCategoryHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	supCatIDStr := c.Param("id")
	supCatID, err := strconv.ParseInt(supCatIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	var req request.UpdateSupportCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err = h.supportCategoryUC.Update(ctx, schema, supCatID, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdateSupportCategorySuccess, "Support category updated successfully", nil)
}

func (h *SupportCategoryHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	supCatIDStr := c.Param("id")
	supCatID, err := strconv.ParseInt(supCatIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	if err := h.supportCategoryUC.Delete(ctx, schema, supCatID); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeDeleteSupportCategorySuccess, "Support category deleted successfully", nil)
}

func (h *SupportCategoryHandler) ListForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	supCatsRes, err := h.supportCategoryUC.ListForGuest(ctx, schema)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"support_categories": supCatsRes,
	})
}
