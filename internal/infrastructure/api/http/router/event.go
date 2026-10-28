package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupEventRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.EventHandler) {
	admin := rg.Group("/admin/events", authMid.IsAuthentication, authMid.HasPermission(model.PermEventManagement))
	{
		admin.POST("", hdl.Create)

		admin.GET("/:id", hdl.GetByID)

		admin.GET("", hdl.ListForAdmin)

		admin.PUT("/:id", hdl.Update)

		admin.DELETE("/:id", hdl.Delete)

		admin.DELETE("", hdl.DeleteBatch)
	}

	guest := rg.Group("/events", authMid.IsGuest)
	{
		guest.GET("", hdl.ListForGuest)

		guest.GET("/:slug", hdl.GetBySlug)
	}
}
