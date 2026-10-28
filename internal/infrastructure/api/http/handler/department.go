package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	departmentUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/department"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	httpRes "github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	departmentUC departmentUC.DepartmentUseCase
}

func NewDepartmentHandler(departmentUC departmentUC.DepartmentUseCase) *DepartmentHandler {
	return &DepartmentHandler{departmentUC}
}

func (h *DepartmentHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.departmentUC.Create(ctx, schema, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	httpRes.APIResponse(c, http.StatusCreated, constants.CodeCreateDepartmentSuccess, "Department created successfully", nil)
}

func (h *DepartmentHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	view := c.Query("view")

	if view == "basic" {
		deptsRes, err := h.departmentUC.ListBasic(ctx, schema)
		if err != nil {
			_ = c.Error(err)
			return
		}

		httpRes.OKResponse(c, gin.H{
			"departments": deptsRes,
		})
		return
	}

	var query request.DepartmentPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	deptsRes, meta, err := h.departmentUC.ListFull(ctx, schema, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	httpRes.OKResponse(c, gin.H{
		"departments": deptsRes,
		"meta":        meta,
	})
}

func (h *DepartmentHandler) GetByIdentifier(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)
	identifier := c.Param("identifier")

	var (
		deptRes *response.BasicDepartmentResponse
		err     error
	)

	if deptID, parseErr := strconv.ParseInt(identifier, 10, 64); parseErr == nil {
		deptRes, err = h.departmentUC.GetByID(ctx, schema, deptID)
	} else {
		deptRes, err = h.departmentUC.GetByName(ctx, schema, identifier)
	}

	if err != nil {
		_ = c.Error(err)
		return
	}

	httpRes.OKResponse(c, gin.H{
		"department": deptRes,
	})
}

func (h *DepartmentHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	deptIDStr := c.Param("id")
	deptID, err := strconv.ParseInt(deptIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	var req request.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.departmentUC.Update(ctx, schema, deptID, userID, req); err != nil {
		_ = c.Error(err)
		return
	}

	httpRes.APIResponse(c, http.StatusOK, constants.CodeUpdateDepartmentSuccess, "Department updated successfully", nil)
}

func (h *DepartmentHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	deptIDStr := c.Param("id")
	deptID, err := strconv.ParseInt(deptIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	if err := h.departmentUC.Delete(ctx, schema, deptID); err != nil {
		_ = c.Error(err)
		return
	}

	httpRes.APIResponse(c, http.StatusOK, constants.CodeDeleteDepartmentSuccess, "Department deleted successfully", nil)
}
