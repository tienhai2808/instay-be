package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupServiceRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.ServiceHandler) {
	admin := rg.Group("/admin/services", authMid.IsAuthentication, authMid.HasRole(model.RoleAdmin))
	{
		admin.POST("", hdl.Create)

		admin.GET("", hdl.ListForAdmin)

		admin.GET("/:id", hdl.GetByID)

		admin.PUT("/:id", hdl.Update)

		admin.DELETE("/:id", hdl.Delete)

		admin.DELETE("", hdl.DeleteBatch)
	}

	guests := rg.Group("/services", authMid.IsGuest)
	{
		guests.GET("", hdl.ListForGuest)

		guests.GET("/:slug", hdl.GetBySlug)
	}
}
