package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupNotificationRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.NotificationHandler) {
	guest := rg.Group("/notifications", authMid.IsGuest)
	{
		guest.GET("/count", hdl.CountUnreadForGuest)

		guest.GET("", hdl.ListForGuest)
	}

	admin := rg.Group("/admin/notifications", authMid.IsAuthentication, authMid.HasRole(model.RoleStaff))
	{
		admin.GET("/count", hdl.CountUnreadForAdmin)

		admin.GET("", hdl.ListForAdmin)
	}
}
