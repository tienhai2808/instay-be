package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupOrderServiceRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.OrderServiceHandler) {
	guest := rg.Group("/order-services", authMid.IsGuest)
	{
		guest.POST("", hdl.CreateForGuest)

		guest.GET("", hdl.ListForGuest)

		guest.PUT("/:id/status", hdl.UpdateStatusForGuest)
	}

	admin := rg.Group("/admin/order-services", authMid.IsAuthentication, authMid.HasPermission(model.PermOrderServiceManagement))
	{
		admin.POST("", hdl.CreateForAdmin)

		admin.GET("", hdl.ListForAdmin)

		admin.GET("/:id", hdl.GetByIDForAdmin)

		admin.PUT("/:id/status", hdl.UpdateStatusForAdmin)
	}
}
