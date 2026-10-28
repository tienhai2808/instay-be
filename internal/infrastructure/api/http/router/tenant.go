package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupTenantRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.TenantHandler) {
	tenant := rg.Group("")
	{
		tenant.GET("", hdl.GetBySchemaName)

		tenant.GET("/settings", hdl.GetSettingsBySchemaName)
		tenant.PUT("/settings", authMid.IsAuthentication, authMid.HasRole(model.RoleAdmin), hdl.UpdateSettings)

		tenant.PUT("", authMid.IsAuthentication, authMid.HasRole(model.RoleAdmin), hdl.Update)
	}
}
