package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	userUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/user"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUC userUC.UserUseCase
}

func NewUserHandler(userUC userUC.UserUseCase) *UserHandler {
	return &UserHandler{userUC}
}

func (h *UserHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userID := c.GetInt64(middleware.CtxUserID)
	if userID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if req.Role == model.RoleAdmin {
		if req.DepartmentID != nil {
			_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
				"field": "departmentid",
				"tag":   "notrequired",
				"param": "",
			}))
			return
		}
	} else {
		if req.DepartmentID == nil {
			_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
				"field": "departmentid",
				"tag":   "required",
				"param": "",
			}))
			return
		}
	}

	id, err := h.userUC.Create(ctx, schema, userID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusCreated, constants.CodeCreateUserSuccess, "User created successfully", gin.H{
		"id": strconv.FormatInt(id, 10),
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	userRes, err := h.userUC.GetByID(ctx, schema, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"user": userRes,
	})
}

func (h *UserHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	var query request.UserPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	usersRes, meta, err := h.userUC.List(ctx, schema, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"users": usersRes,
		"meta":  meta,
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	currentUserID := c.GetInt64(middleware.CtxUserID)
	if currentUserID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if req.Role == model.RoleAdmin {
		if req.DepartmentID != nil {
			_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
				"field": "departmentid",
				"tag":   "notrequired",
				"param": "",
			}))
			return
		}
	} else {
		if req.DepartmentID == nil {
			_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
				"field": "departmentid",
				"tag":   "required",
				"param": "",
			}))
			return
		}
	}

	if err := h.userUC.Update(ctx, schema, userID, currentUserID, req); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdateUserSuccess, "User updated successfully", nil)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	currentUserID := c.GetInt64(middleware.CtxUserID)
	if currentUserID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.UpdateUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	if err := h.userUC.UpdatePassword(ctx, schema, userID, currentUserID, req.NewPassword); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeUpdateUserSuccess, "User password updated successfully", nil)
}

func (h *UserHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		_ = c.Error(exceptions.ErrInvalidID)
		return
	}

	currentUserID := c.GetInt64(middleware.CtxUserID)
	if currentUserID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	if err := h.userUC.Delete(ctx, schema, userID, currentUserID); err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeDeleteUserSuccess, "User deleted successfully", nil)
}

func (h *UserHandler) DeleteBatch(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	currentUserID := c.GetInt64(middleware.CtxUserID)
	if currentUserID == 0 {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	var req request.DeleteBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		field, tag, param := validator.HandleRequestError(err)
		_ = c.Error(exceptions.ErrBadRequest.WithData(gin.H{
			"field": field,
			"tag":   tag,
			"param": param,
		}))
		return
	}

	rowsDeleted, err := h.userUC.DeleteBatch(ctx, schema, currentUserID, req.IDs)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.APIResponse(c, http.StatusOK, constants.CodeDeleteUsersSuccess, "Users deleted successfully", gin.H{
		"count": rowsDeleted,
	})
}
