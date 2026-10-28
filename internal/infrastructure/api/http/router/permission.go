package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupPermissionRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.PermissionHandler) {
	perm := rg.Group("/permissions", authMid.IsAuthentication, authMid.HasRole(model.RoleAdmin))
	{
		perm.GET("", hdl.List)
	}
}
