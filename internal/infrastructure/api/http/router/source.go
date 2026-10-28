package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupSourceRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.SourceHandler) {
	source := rg.Group("/sources", authMid.IsAuthentication)
	{
		source.POST("", authMid.HasRole(model.RoleAdmin), hdl.Create)

		source.PUT("/:id", authMid.HasRole(model.RoleAdmin), hdl.Update)

		source.DELETE("/:id", authMid.HasRole(model.RoleAdmin), hdl.Delete)

		source.GET("", authMid.HasPermission(model.PermBookingManagement), hdl.List)
	}
}
