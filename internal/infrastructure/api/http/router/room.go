package router

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func (r *Router) setupRoomRoutes(rg *gin.RouterGroup, authMid *middleware.AuthMiddleware, hdl *handler.RoomHandler) {
	room := rg.Group("/rooms", authMid.IsAuthentication)
	{
		room.POST("", authMid.HasRole(model.RoleAdmin), hdl.Create)

		room.GET("", authMid.HasRole(model.RoleAdmin), hdl.List)

		room.GET("/timeline", authMid.HasPermission(model.PermBookingManagement), hdl.Timeline)

		room.PUT("/:id", authMid.HasRole(model.RoleAdmin), hdl.Update)

		room.DELETE("/:id", authMid.HasRole(model.RoleAdmin), hdl.Delete)

		room.DELETE("", authMid.HasRole(model.RoleAdmin), hdl.DeleteBatch)
	}
}
