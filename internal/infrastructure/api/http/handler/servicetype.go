package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	serviceTypeUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/servicetype"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type ServiceTypeHandler struct {
	serviceTypeUC serviceTypeUC.ServiceTypeUseCase
}

func NewServiceTypeHandler(serviceTypeUC serviceTypeUC.ServiceTypeUseCase) *ServiceTypeHandler {
	return &ServiceTypeHandler{serviceTypeUC}
}

func (h *ServiceTypeHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.CreateServiceTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.serviceTypeUC.Create(ctx, schema, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusCreated, constants.CodeCreateServiceTypeSuccess, "Service type created successfully", nil)
}

func (h *ServiceTypeHandler) ListForAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	view := c.Query("view")

	if view == "basic" {
		serviceTypesRes, err := h.serviceTypeUC.ListBasicForAdmin(ctx, schema)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response.OKResponse(c, gin.H{
			"service_types": serviceTypesRes,
		})
		return
	}

	var query request.ServiceTypePaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	serviceTypesRes, meta, err := h.serviceTypeUC.ListFullForAdmin(ctx, schema, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"service_types": serviceTypesRes,
		"meta":          meta,
	})
}

func (h *ServiceTypeHandler) ListWithActiveServices(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var query request.ServiceTypeWithActiveServicesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	serviceTypesRes, err := h.serviceTypeUC.ListWithActiveServices(ctx, schema, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"service_types": serviceTypesRes,
	})
}

func (h *ServiceTypeHandler) ListServicesBySlug(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	svcTypeSlug := c.Param("slug")
	include := c.Query("include")
	includeServiceType := include == "service_type"
	var branchID int64
	if branchIDStr := c.Query("branch_id"); branchIDStr != "" {
		var err error
		branchID, err = strconv.ParseInt(branchIDStr, 10, 64)
		if err != nil {
			_ = c.Error(exceptions.ErrInvalidID)
			return
		}
	}

	svcTypeRes, svcsRes, err := h.serviceTypeUC.ListServicesBySlug(ctx, schema, svcTypeSlug, includeServiceType, branchID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res := gin.H{
		"services": svcsRes,
	}

	if includeServiceType {
		res["service_type"] = svcTypeRes
	}

	response.OKResponse(c, res)
}

func (h *ServiceTypeHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	svcTypeIDStr := c.Param("id")
	svcTypeID, err := strconv.ParseInt(svcTypeIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	var req request.UpdateServiceTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.serviceTypeUC.Update(ctx, schema, svcTypeID, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdateServiceTypeSuccess, "Service type updated successfully", nil)
}

func (h *ServiceTypeHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	svcTypeIDStr := c.Param("id")
	svcTypeID, err := strconv.ParseInt(svcTypeIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	if err := h.serviceTypeUC.Delete(ctx, schema, svcTypeID); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeDeleteServiceTypeSuccess, "Service type deleted successfully", nil)
}

func (h *ServiceTypeHandler) ListForGuest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	serviceTypesRes, err := h.serviceTypeUC.ListForGuest(ctx, schema)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"service_types": serviceTypesRes,
	})
}
