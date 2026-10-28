package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupDepartmentRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.DepartmentHandler) {
	admin := rg.Group("/admin/departments", authMid.IsAuthentication, authMid.HasRole(model.RoleAdmin))
	{
		admin.POST("", hdl.Create)

		admin.GET("", hdl.List)

		admin.PUT("/:id", hdl.Update)

		admin.DELETE("/:id", hdl.Delete)
	}

	guest := rg.Group("/departments", authMid.IsGuest)
	{
		guest.GET("/:identifier", hdl.GetByIdentifier)
	}
}
