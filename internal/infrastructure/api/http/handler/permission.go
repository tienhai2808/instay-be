package handler

import (
	"context"
	"time"

	permissionUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/permission"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionUC permissionUC.PermissionUseCase
}

func NewPermissionHandler(permissionUC permissionUC.PermissionUseCase) *PermissionHandler {
	return &PermissionHandler{permissionUC}
}

func (h *PermissionHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	schema := c.GetString(middleware.CtxSchema)

	permsRes, err := h.permissionUC.List(ctx, schema)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.OKResponse(c, gin.H{
		"permissions": permsRes,
	})
}
