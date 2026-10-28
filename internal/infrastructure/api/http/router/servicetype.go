package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupServiceTypeRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.ServiceTypeHandler) {
	admin := rg.Group("/admin/service-types", authMid.IsAuthentication)
	{
		admin.POST("", authMid.HasRole(model.RoleAdmin), hdl.Create)

		admin.GET("", authMid.HasRole(model.RoleAdmin), hdl.ListForAdmin)

		admin.GET("/services", authMid.HasPermission(model.PermBookingManagement), hdl.ListWithActiveServices)

		admin.PUT("/:id", authMid.HasRole(model.RoleAdmin), hdl.Update)

		admin.DELETE("/:id", authMid.HasRole(model.RoleAdmin), hdl.Delete)
	}

	guest := rg.Group("/service-types")
	{
		guest.GET("", hdl.ListForGuest)

		guest.GET("/:slug/services", hdl.ListServicesBySlug)
	}
}
