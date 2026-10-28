package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupSupportCategoryRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.SupportCategoryHandler) {
	admin := rg.Group("/admin/support-categories", authMid.IsAuthentication, authMid.HasRole(model.RoleAdmin))
	{
		admin.POST("", hdl.Create)

		admin.PUT("/:id", hdl.Update)

		admin.GET("", hdl.ListForAdmin)

		admin.DELETE("/:id", hdl.Delete)
	}

	guest := rg.Group("/support-categories", authMid.IsGuest)
	{
		guest.GET("", hdl.ListForGuest)
	}
}
