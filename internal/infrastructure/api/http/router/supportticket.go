package router

import (
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupSupportTicketRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.SupportTicketHandler) {
	guest := rg.Group("/support-tickets", authMid.IsGuest)
	{
		guest.POST("", hdl.Create)

		guest.GET("", hdl.ListForGuest)

		guest.PUT("/:id/status", hdl.UpdateStatusForGuest)
	}

	admin := rg.Group("/admin/support-tickets", authMid.IsAuthentication)
	{
		admin.GET("", hdl.ListForAdmin)

		admin.PUT("/:id/status", hdl.UpdateStatusForAdmin)
	}
}
