package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupFloorRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.FloorHandler) {
	floor := rg.Group("/floors", authMid.IsAuthentication)
	{
		floor.POST("", authMid.HasRole(model.RoleAdmin), hdl.Create)

		floor.GET("", authMid.HasPermission(model.PermBookingManagement), hdl.List)

		floor.PUT(":id", authMid.HasRole(model.RoleAdmin), hdl.Update)

		floor.DELETE("id", authMid.HasRole(model.RoleAdmin), hdl.Delete)
	}
}
