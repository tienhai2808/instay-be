package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupCustomerRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.CustomerHandler) {
	guest := rg.Group("/customers", authMid.IsGuest)
	{
		guest.POST("", hdl.Create)
	}

	admin := rg.Group("/admin/customers", authMid.IsAuthentication, authMid.HasRole(model.RoleAdmin))
	{
		admin.GET("", hdl.List)

		admin.POST("/export", hdl.Export)
	}
}
